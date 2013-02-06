package coremidi

type PacketWriter struct {
	port        OutputPort
	destination Destination
}

func NewPacketWriter(port OutputPort, destination Destination) *PacketWriter {
	return &PacketWriter{port, destination}
}

func (writer *PacketWriter) Write(p []byte) (n int, err error) {
	packet := NewPacket(p)
	err = packet.Send(&writer.port, &writer.destination)

	n = len(p)

	return
}
