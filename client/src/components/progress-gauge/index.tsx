// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates 
// SPDX-License-Identifier: MIT 

import React, { useState, useEffect, useMemo, useRef } from 'react'
import { View, Text } from '@tarojs/components'

import './index.less'

interface ProgressGaugeProps {
  total: number
  current: number
}

const ProgressGauge: React.FC<ProgressGaugeProps> = ({
  total,
  current
}) => {
  const [animatedPercentage, setAnimatedPercentage] = useState(0)
  const [isVisible, setIsVisible] = useState(false)
  const animationRef = useRef<number | null>(null)
  const startTimeRef = useRef<number | null>(null)

  const percentage = useMemo(() => {
    return Math.round((current / total) * 100)
  }, [current, total])

  const dynamicContent = useMemo(() => {
    if (percentage >= 80) {
      return {
        title: `已完成 ${percentage}%`,
        subtitle: '语言的底气，就藏在你记下的每一个单词里',
        gradientColors: ['#00CC6A', '#6BF18C'] // Green gradient
      }
    } else if (percentage >= 50) {
      return {
        title: `已完成 ${percentage}%`,
        subtitle: '语言的底气，就藏在你记下的每一个单词里',
        gradientColors: ['#0578FF', '#CBE8FF'] // Blue gradient
      }
    } else {
      return {
        title: `已完成 ${percentage}%`,
        subtitle: '积累的力量是惊人的，继续加油吧！',
        gradientColors: ['#FF9C1A', '#FFD836'] // Orange gradient
      }
    }
  }, [percentage])

  const easeOutCubic = (t: number): number => {
    return 1 - Math.pow(1 - t, 3)
  }

  const animateProgress = () => {
    const duration = 2000

    const animate = (currentTime: number) => {
      if (startTimeRef.current === null) {
        startTimeRef.current = currentTime
      }

      const elapsed = currentTime - startTimeRef.current
      const progress = Math.min(elapsed / duration, 1)

      // Apply easing function
      const easedProgress = easeOutCubic(progress)
      const currentValue = easedProgress * percentage

      setAnimatedPercentage(currentValue)

      if (progress < 1) {
        animationRef.current = requestAnimationFrame(animate)
      } else {
        // Animation complete, ensure final value is accurate
        setAnimatedPercentage(percentage)
        animationRef.current = null
        startTimeRef.current = null
      }
    }

    animationRef.current = requestAnimationFrame(animate)
  }

  useEffect(() => {
    // Clean up previous animation
    if (animationRef.current) {
      cancelAnimationFrame(animationRef.current)
      animationRef.current = null
      startTimeRef.current = null
    }

    // Reset animation state
    setAnimatedPercentage(0)
    setIsVisible(false)

    const timer = setTimeout(() => {
      setIsVisible(true)
      animateProgress()
    }, 200) // Slight delay to allow page rendering to complete

    return () => {
      clearTimeout(timer)
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current)
      }
    }
  }, [percentage])

  // Clean up animation when component unmounts
  useEffect(() => {
    return () => {
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current)
      }
    }
  }, [])

  // Calculate size proportions
  const svgWidth = 226;
  const svgHeight = 130;
  // const svgWidth = width
  // const svgHeight = width * 0.6 // Height is 60% of width
  const centerX = svgWidth / 2
  // const centerY = svgHeight * 0.85 // Center position
  const centerY = svgHeight - 10
  const radius = svgWidth * 0.46 // Radius
  // const strokeWidth = radius * 0.08 // Progress bar width

  // Calculate progress
  const circumference = Math.PI * radius
  const progressLength = (animatedPercentage / 100) * circumference

  // Generate semi-circle path
  const createSemiCirclePath = (r: number) => {
    return `M ${centerX - r} ${centerY} A ${r} ${r} 0 0 1 ${centerX + r} ${centerY}`
  }
  const gradientColors = dynamicContent.gradientColors // Get gradient colors from dynamicContent

  return (
    <View className={`progress-gauge__main-container ${isVisible ? 'progress-gauge__visible' : ''}`}>
      <Text className='progress-gauge__percent-start'>0</Text>
      <Text className='progress-gauge__percent-sign'>%</Text>
      <View className='progress-gauge__middle'>
        <svg viewBox={`0 0 ${svgWidth} ${svgHeight}`} className='progress-gauge__svg' preserveAspectRatio='xMidYMid meet'>
          <defs>
            <linearGradient id='progressGradient' x1='0%' y1='0%' x2='100%' y2='0%'>
              <stop offset='0%' stopColor={gradientColors[0]} />
              <stop offset='100%' stopColor={gradientColors[1]} />
            </linearGradient>
          </defs>
          <path
            d={createSemiCirclePath(radius)}
            fill='none'
            stroke='#32374A0D'
            strokeWidth='4'
            strokeLinecap='round'
          />
          <path
            d={createSemiCirclePath(radius)}
            fill='none'
            stroke='url(#progressGradient)'
            strokeWidth='12'
            strokeLinecap='round'
            strokeDasharray={`${progressLength} ${circumference}`}
            className='progress-gauge__progress-path'
            style={{
              transition: isVisible ? 'none' : 'stroke-dasharray 0.3s ease'
            }}
          />
        </svg>
        <View className='progress-gauge__todays_progress'>
          <Text className='progress-gauge__todays_progress__label'>今日已学完</Text>
          <Text className='progress-gauge__todays_progress__finished'>{current}</Text>
          <Text className='progress-gauge__todays_progress__target'>今日需学 {total}</Text>
        </View>
      </View>
      <Text className='progress-gauge__percent-start'>100</Text>
      <Text className='progress-gauge__percent-sign'>%</Text>
    </View>
  )
}

export default ProgressGauge
