package util

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GerarCpf() string {
	num1 := aleatorio()
	num2 := aleatorio()
	num3 := aleatorio()

	dig1 := dig(num1, num2, num3, "")
	dig2 := dig(num1, num2, num3, dig1)

	return fmt.Sprint(num1, num2, num3, dig1, dig2)
}

func dig(n1, n2, n3, n4 string) string {
	var nums []string
	nums = append(nums, strings.Split(n1, "")...)
	nums = append(nums, strings.Split(n2, "")...)
	nums = append(nums, strings.Split(n3, "")...)

	if n4 == "" {
		nums = append(nums, "0")
	} else {
		nums = append(nums, n4, "")
	}

	x := 0
	i := len(nums)
	j := 0

	for i >= 2 {
		n, err := strconv.Atoi(nums[j])

		if err != nil {
			panic(err)
		}

		x += n * i

		i--
		j++
	}

	y := x % 11

	if y < 2 {
		return "0"
	}

	return fmt.Sprint(11-y)
}

func aleatorio() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	i := r.Intn(999)

	return fmt.Sprintf("%03d", i)
}
