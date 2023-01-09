package server

func Multiplex(req map[string][100]float32) map[string][100]float32 {
	resp := make(map[string][100]float32)
	for mId := range req {
		go GetMemberChannel(mId, req, &resp[mId])
	}
	return resp
}

func GetMemberChannel(memberId string, req map[string][100]float32, returnArr *[100]float32) {
	var resp [100]float32
	for reqMemberId := range req {
		if reqMemberId == memberId {
			continue
		}
		resp = Sum(resp, req[reqMemberId])
	}
	returnArr = &resp
}

func Sum(first [100]float32, second [100]float32) [100]float32 {
	var res [100]float32
	for i := 0; i < 100; i++ {
		res[i] = first[i] + second[i]
	}
	return res
}
