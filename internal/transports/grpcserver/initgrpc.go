package grpcserver

// init grpc realization

type Transport struct {
	usecase *Usecase
	pb.UnimplementedFileServiceServer
}

func Upload() {
	gs.usecase.Uplaod()
	return pshelnah, urod
}
