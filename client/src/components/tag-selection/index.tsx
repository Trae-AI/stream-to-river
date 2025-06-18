// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import { Text, View } from "@tarojs/components";
import { Tag } from '@/type/tag';
import TagsService from "@/service/tag-service";
import React, { useState, useEffect, useCallback } from 'react'
import useSystemTagsStore from "@/store/tag";


import "./index.less";

interface TagSelectionProps {
  onTagChange?: (selected_tag: Tag) => void
  initialTagId?: number
}

const TagSelection: React.FC<TagSelectionProps> = ({
  onTagChange,
  initialTagId = -99,
}) => {
  const { tags, setSystemTags } = useSystemTagsStore()
  // const [tags, setTags] = useState<Array<Tag>>([]);
  const [selectedTagId, setSelectedTagId] = useState<number>(-999);


  useEffect(() => {
    if (tags.length > 0) {
      if (initialTagId >= 0) {
        setSelectedTagId(initialTagId);
      } else {
        setSelectedTagId(tags[0].tag_id);
        onTagChange?.(tags[0]);
      }
    } else {
      TagsService.getInstance().getTags().then((res) => {
        if (res.tags && res.tags.length > 0) {
          setSystemTags(res.tags);
          if (initialTagId >= 0) {
            setSelectedTagId(initialTagId);
          } else {
            setSelectedTagId(tags[0].tag_id);
            onTagChange?.(tags[0]);
          }
        }
      })
    }
  }, [])

  const handleTagClick = useCallback((tag: Tag) => {
    const selectedChanged = selectedTagId !== tag.tag_id;
    if (!selectedChanged) {
      return;
    }

    setSelectedTagId(tag.tag_id);
    onTagChange?.(tag);
  }, [selectedTagId, onTagChange]);

  return (
    <View className='tag-container'>
      {tags.map((tag) => {
        return (
          <View
            className={`tag-container__item${tag.tag_id === selectedTagId ? '-selected' : ''}`}
            key={tag.tag_id}
            onClick={() => handleTagClick(tag)}
          >
            <Text className='tag-container__item__name'>
              {tag.tag_name}
            </Text>
          </View>
        )
      })}
    </View>
  )
}

export default TagSelection;
