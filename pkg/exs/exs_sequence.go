package exs

import "k8s.io/klog"

func (exs *EXS) ReadSequences() error {

	for i := int32(0); i < int32(len(exs.Groups)); i++ {
		group := exs.Groups[i]
		klog.V(2).Infof("group %d: %s selectGroup: %d\n", i, group.Name, group.SelectGroup)
		//if group.SelectGroup == -1 {
		//	continue
		//}
		foundInSequence := false
		for _, sequence := range exs.Sequences {
			if contains(sequence, i) {
				foundInSequence = true
				klog.V(2).Infof("found in sequence %d: %+v\n", i, sequence)
				break
			}
		}
		if foundInSequence {
			continue
		}

		gid := i
		sequence := []int32{}
		cont := true
		for cont {
			cont = false
			for j := int32(0); j < int32(len(exs.Groups)); j++ {
				g := exs.Groups[j]
				if g.SelectGroup == gid && (j != g.SelectGroup) && !contains(sequence, gid) {
					klog.V(2).Infof("adding %d to sequence\n", gid)
					sequence = append(sequence, gid)
					gid = j
					cont = true
					break
				}
			}
		}
		// now that we're at the start of the chain, simply follow it to the end
		sequence = []int32{}
		for gid != int32(-1) && !contains(sequence, gid) {
			klog.V(2).Infof("gid: %d\n", gid)
			sequence = append([]int32{gid}, sequence...)
			gid = exs.Groups[gid].SelectGroup
		}
		if len(sequence) > 1 {
			exs.Sequences = append(exs.Sequences, sequence)
		}
	}

	return nil
}

// convertSeqNumbers converts the sequence numbers to the actual sequence number
func (exs *EXS) ConvertSeqNumbers() error {
	for i := int32(0); i < int32(len(exs.Groups)); i++ {
		group := exs.Groups[i]
		group.SelectNumber = 0
		for j := int32(0); j < int32(len(exs.Sequences)); j++ {
			sequence := exs.Sequences[j]
			if contains(sequence, i) {
				group.SelectNumber = uint8(indexOf(sequence, j) + 1)
			}
		}
	}
	return nil
}

func contains(slice []int32, item int32) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

func indexOf(slice []int32, item int32) int32 {
	for i, v := range slice {
		if v == item {
			return int32(i)
		}
	}
	return -1
}
