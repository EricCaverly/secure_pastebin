package main

import (
	"fmt"
	"strconv"
	"strings"
)

func check_valid_ranges(ranges string) error {
	_, err := within_ranges("0.0.0.0", ranges)
	return err
}

func ipa_to_ipi(ip string) (uint32, error) {
	segments := strings.Split(ip, ".")

	if len(segments) != 4 {
		return 0, fmt.Errorf("invalid number of segments in ip address")
	}

	i_segs := [4]int{}

	var err error

	for i := range 4 {
		i_segs[i], err = strconv.Atoi(segments[i])
		if err != nil {
			return 0, err
		}

		if i_segs[i] > 255 || i_segs[i] < 0 {
			return 0, fmt.Errorf("invalid %d segment number", i)
		}
	}

	ip_i := i_segs[0]<<24 + i_segs[1]<<16 + i_segs[2]<<8 + i_segs[3]

	return uint32(ip_i), nil
}

func within_ranges(ip string, ranges string) (bool, error) {
	ip_i, err := ipa_to_ipi(ip)
	if err != nil {
		return false, fmt.Errorf("client IP is invalid")
	}

	ranges_a := strings.SplitSeq(ranges, ",")

	for rr := range ranges_a {

		segments := strings.Split(strings.TrimSpace(rr), "/")
		if len(segments) == 1 {
			segments = append(segments, "32")
		}

		mask_n, err := strconv.Atoi(segments[1])
		if err != nil || mask_n < 0 || mask_n > 32 {
			return false, fmt.Errorf("invalid subnet mask")
		}

		mask := ^uint32((1 << (32 - mask_n)) - 1)
		// fmt.Printf("%032b\n", mask)

		range_ip, err := ipa_to_ipi(segments[0])
		if err != nil {
			return false, err
		}

		// fmt.Printf("%032b\n%032b\n", (range_ip & mask), (ip_i & mask))

		if (range_ip & mask) == (ip_i & mask) {
			return true, nil
		}
	}

	return false, nil
}
