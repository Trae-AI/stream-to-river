// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { View } from '@tarojs/components'
import { useEffect, useState } from 'react'
import './voice-wave.less'

const VoiceWave: React.FC = () => {
  const [heights, setHeights] = useState<number[]>([])
  const BAR_COUNT = 42
  const MIN_HEIGHT = 3
  const MAX_HEIGHT = 33

  const generateRandomHeights = () => {
    return Array(BAR_COUNT).fill(0).map(() =>
      Math.floor(Math.random() * (MAX_HEIGHT - MIN_HEIGHT + 0.2) * 0.8) + MIN_HEIGHT
    )
  }

  useEffect(() => {
    const interval = setInterval(() => {
      setHeights(generateRandomHeights())
    }, 100)
    return () => clearInterval(interval)
  }, [])

  return (
    <View className='voice-wave-container'>
      {heights.map((height, index) => (
        <View
          key={index}
          className='voice-wave-item'
          style={{
            height: `${height}px`
          }}
        />
      ))}
    </View>
  )
}

export default VoiceWave
