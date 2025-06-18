// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import React, { useState, useEffect } from "react";
import { View, Text } from "@tarojs/components";
import "./index.less";

interface CategoryCardProps {
  /**
   * card title
   */
  title: string;
  /**
   * icon path
   */
  iconSrc: React.FC<React.SVGProps<SVGSVGElement>>;
  /**
   * description
   */
  description: string;
  /**
   * click event
   */
  onClick?: () => void;
  /**
   * class name
   */
  className?: string;
}


const CategoryCard: React.FC<CategoryCardProps> = ({
  title,
  iconSrc,
  description,
  onClick,
  className = "",
}) => {
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    const timer = setTimeout(() => {
      setIsVisible(true);
    }, 100);
    return () => clearTimeout(timer);
  }, []);
  const IconSrc = iconSrc;
  return (
    <View
      className={`category-card ${
        isVisible ? "category-card--visible" : ""
      } ${className}`}
      onClick={onClick}
    >
      <IconSrc className="category-card__icon" />
      <View className="category-card__content">
        <Text className="category-card__title">{title}</Text>
        <Text className="category-card__description">{description}</Text>
      </View>
    </View>
  );
};

export default CategoryCard;
