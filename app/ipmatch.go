/*********************************
 *  File     : ipmatch.go
 *  Purpose  : Helper functions for IP range enclusion checking
 *  Authors  : Eric Caverly
 */

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

	// Ensure we have a & b & c & d
	if len(segments) != 4 {
		return 0, fmt.Errorf("invalid number of segments in ip address")
	}

	// Convert a, b, c, d into integers and ensure they are between 0 and 255
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

	// Combine a + b + c + d into an unsigned 32 bit integer representation of a.b.c.d
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

		// Split a.b.c.d/m into a.b.c.d and m
		segments := strings.Split(strings.TrimSpace(rr), "/")
		if len(segments) == 1 {
			segments = append(segments, "32")
		}

		// Turn m into an int. Verify it's between 0 and 32 inclusive
		mask_n, err := strconv.Atoi(segments[1])
		if err != nil || mask_n < 0 || mask_n > 32 {
			return false, fmt.Errorf("invalid subnet mask")
		}

		// Transform m into the binary subnet mask representation of the simplified network prefix
		mask := ^uint32((1 << (32 - mask_n)) - 1)
		// fmt.Printf("%032b\n", mask)

		range_ip, err := ipa_to_ipi(segments[0])
		if err != nil {
			return false, err
		}

		// fmt.Printf("%032b\n%032b\n", (range_ip & mask), (ip_i & mask))

		// Compare the network portions of the IP in question and this range IP
		// If they match, then the IP in question is within the range IP
		if (range_ip & mask) == (ip_i & mask) {
			return true, nil
		}
	}

	return false, nil
}
