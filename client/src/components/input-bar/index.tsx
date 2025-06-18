// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { useMemo, useRef, useState } from "react";
import { ITouchEvent, Text, Textarea, View } from "@tarojs/components";
import "./index.less";
import { ReactComponent as Mic } from "@/assets/icons/mic.svg";
import { ReactComponent as Camera } from "@/assets/icons/camera.svg";
import { ReactComponent as SendIcon } from "@/assets/icons/send.svg";
import { ReactComponent as SendDisabledIcon } from "@/assets/icons/send-disabled.svg";
import { ReactComponent as SendGeneratingIcon } from "@/assets/icons/send-generating.svg";
import VoiceWave from "./voice-wave";
import { ReactComponent as KeyboardIcon } from "@/assets/icons/keyboard.svg";
import { InputBarStatus } from "./types";
import { LONG_PRESS_TIME_MS } from "./consts";
import { useAudioRecording } from "@/hooks/asr";
import Taro from "@tarojs/taro";
import { sleep } from "@/pages/index/utils";

enum SendButtonStatus {
  Disabled = 0,
  Enable = 1,
  Generating = 2,
}

const MicButton = ({ switchToAudio }: { switchToAudio: () => void }) => {
  return <Mic className="input-bar-button" onClick={switchToAudio}></Mic>;
};
const CameraButton = ({ switchToCamera }: { switchToCamera: () => void }) => {
  return (
    <Camera className="input-bar-button" onClick={switchToCamera}></Camera>
  );
};
const SendButton = ({
  status,
  onSend,
  onCancel,
}: {
  status: SendButtonStatus;
  onSend: () => void;
  onCancel: () => void;
}) => {
  if (status === SendButtonStatus.Enable) {
    return <SendIcon className="input-bar-button-send" onClick={onSend} />;
  } else if (status === SendButtonStatus.Generating) {
    return (
      <SendGeneratingIcon
        className="input-bar-button-send"
        onClick={onCancel}
      />
    );
  } else {
    return <SendDisabledIcon className="input-bar-button-send" />;
  }
};
const KeyboardButton = ({
  switchToKeyboard,
}: {
  switchToKeyboard: () => void;
}) => {
  return (
    <KeyboardIcon className="input-bar-button" onClick={switchToKeyboard} />
  );
};

const Init = ({
  setInputBarStatus,
  switchToCamera,
}: {
  setInputBarStatus: (status: InputBarStatus) => void;
  switchToCamera: () => void;
}) => {
  return (
    <View className="input-bar-init-container">
      <Text
        className="input-bar-init-placeholder"
        onClick={() => setInputBarStatus(InputBarStatus.Typing)}
      >
        查询翻译或按住说话...
      </Text>
      <View className="gap-16">
        <MicButton
          switchToAudio={() => setInputBarStatus(InputBarStatus.AudioInit)}
        />
        <CameraButton switchToCamera={switchToCamera} />
      </View>
    </View>
  );
};

const Typing = ({
  hasHistoryMessage,
  isGenerating,
  send,
  cancel,
  setInputBarStatus,
  switchToCamera,
}: {
  hasHistoryMessage: boolean;
  isGenerating: boolean;
  send: (text: string) => void;
  cancel: () => void;
  setInputBarStatus: (status: InputBarStatus) => void;
  switchToCamera: () => void;
}) => {
  const [text, setText] = useState<string>("");

  const sendButtonStatus = useMemo(() => {
    if (isGenerating) {
      return SendButtonStatus.Generating;
    }
    if (text.length === 0) {
      return SendButtonStatus.Disabled;
    }
    return SendButtonStatus.Enable;
  }, [isGenerating, text]);

  const handleSend = () => {
    if (text.length === 0) {
      return;
    }
    send(text);
    setText("");
  };

  const placeholder = useMemo(() => {
    if (isGenerating) {
      return "生成中...";
    } else if (hasHistoryMessage) {
      return "继续对话";
    }
    return "请输入";
  }, [hasHistoryMessage, isGenerating]);

  const handleCancel = () => {
    cancel();
    setText("");
  };

  return (
    <View className="input-bar-typing-container">
      <Textarea
        disabled={isGenerating}
        value={text}
        placeholder={placeholder}
        autoFocus
        autoHeight
        className="input-bar-typing-input"
        placeholderClass="input-bar-typing-input-placeholder"
        onInput={(e) => setText(e.detail.value)}
        maxlength={1000}
        style={{ maxHeight: '200px', overflow: 'auto'}}
      />
      <View className="input-bar-typing-bottom">
        <View className="gap-16">
          <MicButton
            switchToAudio={() => setInputBarStatus(InputBarStatus.AudioInit)}
          />
          <CameraButton switchToCamera={switchToCamera} />
        </View>
        <SendButton
          status={sendButtonStatus}
          onSend={handleSend}
          onCancel={handleCancel}
        />
      </View>
    </View>
  );
};

const AudioInit = ({
  switchToCamera,
  setInputBarStatus,
}: {
  switchToCamera: () => void;
  setInputBarStatus: (status: InputBarStatus) => void;
}) => {
  return (
    <View className="input-bar-init-container">
      <CameraButton switchToCamera={switchToCamera} />
      <Text className="audio-init-tips">按住说话</Text>
      <KeyboardButton
        switchToKeyboard={() => setInputBarStatus(InputBarStatus.Typing)}
      />
    </View>
  );
};

enum AudioInputStatus {
  Hidden = 0,
  Speaking = 1,
  Cancel = 2,
}

const AudioInput = ({ status }: { status: AudioInputStatus }) => {
  return (
    <View
      className="audio-container"
      style={{
        backgroundColor:
          status === AudioInputStatus.Cancel ? "#FA3A3A" : "#0079ef",
      }}
    >
      <Text className="audio-tips">
        {status === AudioInputStatus.Cancel
          ? "取消发送"
          : "松手发送语音，上划取消发送"}
      </Text>
      <View className="audio-bar">
        <VoiceWave />
      </View>
    </View>
  );
};

type InputBarProps = {
  hasHistoryMessage: boolean;
  isGenerating: boolean;
  send: (text: string) => void;
  cancel: () => void;
  sendImage: (image: string) => void;
};

const InputBar = (props: InputBarProps) => {
  const [status, setStatus] = useState<InputBarStatus>(InputBarStatus.Init);
  const [audioStatus, setAudioStatus] = useState<AudioInputStatus>(
    AudioInputStatus.Hidden
  );
  const [isLongPress, setIsLongPress] = useState(false);

  const switchToCamera = async () => {
    Taro.showToast({
      title: "敬请期待!",
      icon: "none",
    })
  };

  const touchStartY = useRef(0);
  const timer = useRef<NodeJS.Timeout>();

  const { startRecording, stopRecording, cancelRecording } =
    useAudioRecording();

  const handleTouchStart = (e: ITouchEvent) => {
    if (status !== InputBarStatus.Init && status !== InputBarStatus.AudioInit) {
      return;
    }

    touchStartY.current = e.touches[0].clientY;
    timer.current = setTimeout(async () => {
      try {
        setIsLongPress(true);
        setAudioStatus(AudioInputStatus.Speaking);

        const success = await startRecording();
        if (!success) {
          setIsLongPress(false);
          setAudioStatus(AudioInputStatus.Hidden);
        }
      } catch (error) {
        clearTimeout(timer.current);
        setIsLongPress(false);
        setAudioStatus(AudioInputStatus.Hidden);
      }
    }, LONG_PRESS_TIME_MS);
  };

  const handleTouchMove = (e: ITouchEvent) => {
    if (!isLongPress) return;

    const moveY = e.touches[0].clientY;
    const diff = touchStartY.current - moveY;

    if (diff > 50) {
      setAudioStatus(AudioInputStatus.Cancel);
    } else {
      setAudioStatus(AudioInputStatus.Speaking);
    }
  };

  const handleTouchEnd = async (e: ITouchEvent) => {
    if (timer.current) {
      clearTimeout(timer.current);
    }

    if (isLongPress) {
      setIsLongPress(false);
      if (audioStatus === AudioInputStatus.Speaking) {
        const asyncStopRecording = async () => {
          try {
            await sleep(500);
            const recognizedText = await stopRecording();
            if (recognizedText) {
              // setStatus(InputBarStatus.Typing);
              props.send(recognizedText);
            }
          } catch (error) {
            console.error("Speech recognition failed:", error);
            if (error?.status === 500) {
              Taro.showToast({
                title: "录音时间太短, 请重试",
                icon: "none",
              });
            } else if (error?.data?.code === 1013) {
              Taro.showToast({
                title: "无有效人声语音，请重试",
                icon: "none",
              });
            } else {
              Taro.showToast({
                title: "录音太短或无人声，请尝试再次发起",
                icon: "none",
              });
            }
          }
        };
        asyncStopRecording();
      } else {
        cancelRecording();
      }

      setAudioStatus(AudioInputStatus.Hidden);
      setStatus(InputBarStatus.AudioInit);
    }
  };

  return (
    <View
      className="input-bar-wrapper"
      onTouchStart={handleTouchStart}
      onTouchMove={handleTouchMove}
      onTouchEnd={handleTouchEnd}
      onTouchForceChange={() => {
        alert("xx");
      }}
    >
      {status === InputBarStatus.Init && (
        <Init setInputBarStatus={setStatus} switchToCamera={switchToCamera} />
      )}
      {status === InputBarStatus.Typing && (
        <Typing
          hasHistoryMessage={props.hasHistoryMessage}
          isGenerating={props.isGenerating}
          send={props.send}
          cancel={props.cancel}
          setInputBarStatus={setStatus}
          switchToCamera={switchToCamera}
        />
      )}
      {status === InputBarStatus.AudioInit && (
        <AudioInit
          switchToCamera={switchToCamera}
          setInputBarStatus={setStatus}
        />
      )}
      {audioStatus !== AudioInputStatus.Hidden && (
        <AudioInput status={audioStatus} />
      )}
    </View>
  );
};

export default InputBar;
