package utils

import (
	"gameserver/internal/pb"
)

func ConverterArray2ResInfo(resInfos [][]int) []*pb.ResInfo {
	retRes := make([]*pb.ResInfo, 0, len(resInfos))
	for _, res := range resInfos {
		retRes = append(retRes, &pb.ResInfo{
			BaseType: int32(res[0]),
			Id:       int32(res[1]),
			Count:    int32(res[2]),
		})
	}
	return retRes
}

func ConverterMapResInfo(baseType int, resInfos map[int]int) []*pb.ResInfo {
	retRes := make([]*pb.ResInfo, 0, len(resInfos))
	for k, v := range resInfos {
		retRes = append(retRes, &pb.ResInfo{BaseType: int32(baseType), Id: int32(k), Count: int32(v)})
	}
	return retRes
}
