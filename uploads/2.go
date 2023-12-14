package main

import (
	"fmt"
)

func localAlignment(Pm, T string, matchScore, mismatchPenalty, gapPenalty int) (string, string, int) {
	m := len(Pm)
	n := len(T)

	// 创建动态规划矩阵
	matrix := make([][]int, m+1)
	for i := range matrix {
		matrix[i] = make([]int, n+1)
	}

	// 记录最大得分和位置
	maxScore := 0
	maxI, maxJ := 0, 0

	// 填充动态规划矩阵并找到最大得分
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			match := matrix[i-1][j-1] + matchScore
			delete := matrix[i-1][j] + gapPenalty
			insert := matrix[i][j-1] + gapPenalty
			mismatch := matrix[i-1][j-1] - mismatchPenalty

			matrix[i][j] = Max(0, match, delete, insert, mismatch)

			// 更新最大得分和位置
			if matrix[i][j] > maxScore {
				maxScore = matrix[i][j]
				maxI, maxJ = i, j
			}
		}
	}

	// 回溯路径并获取局部比对
	localAlignmentPm := ""
	localAlignmentT := ""
	i, j := maxI, maxJ
	for i > 0 && j > 0 && matrix[i][j] > 0 {
		currentScore := matrix[i][j]
		match := matrix[i-1][j-1] + matchScore
		delete := matrix[i-1][j] + gapPenalty
		insert := matrix[i][j-1] + gapPenalty
		mismatch := matrix[i-1][j-1] - mismatchPenalty

		if currentScore == match {
			localAlignmentPm = string(Pm[i-1]) + localAlignmentPm
			localAlignmentT = string(T[j-1]) + localAlignmentT
			i--
			j--
		} else if currentScore == delete {
			localAlignmentPm = string(Pm[i-1]) + localAlignmentPm
			localAlignmentT = "-" + localAlignmentT
			i--
		} else if currentScore == insert {
			localAlignmentPm = "-" + localAlignmentPm
			localAlignmentT = string(T[j-1]) + localAlignmentT
			j--
		} else if currentScore == mismatch {
			localAlignmentPm = string(Pm[i-1]) + localAlignmentPm
			localAlignmentT = string(T[j-1]) + localAlignmentT
			i--
			j--
		}
	}

	return localAlignmentPm, localAlignmentT, maxScore

}

func Max(values ...int) int {
	maxValue := values[0]
	for _, value := range values {
		if value > maxValue {
			maxValue = value
		}
	}
	return maxValue
}

func main() {
	Pm := "ABCABC"    // 重复的模式
	T := "BABCABCXYZ" // 文本

	matchScore := 1      // 匹配得分
	mismatchPenalty := 1 // 不匹配惩罚
	gapPenalty := -1     // 间隙惩罚

	seq1, seq2, score := localAlignment(Pm, T, matchScore, mismatchPenalty, gapPenalty)

	fmt.Println("局部比对序列1:", seq1)
	fmt.Println("局部比对序列2:", seq2)
	fmt.Println("最大得分:", score)

}
