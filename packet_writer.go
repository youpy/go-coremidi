package coremidi

type PacketWriter struct {
	port        *OutputPort
	destination *Destination
}

func (writer *PacketWriter) Write(p []byte) (n int, err error) {
	packet := NewPacket(p)
	err = packet.Send(writer.port, writer.destination)

	n = len(p)

	return
}
