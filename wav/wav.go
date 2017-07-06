package wav
//I ONLY HAD TO MAKE Create() RETURN A *Wav BECAUSE Wav.Patch() WOULDN'T UPDATE THE Wav ITSELF, IT NEEDED TO BE *Wav.Patch()

import (
	"fmt"
	"encoding/binary" //for converting ints to []bytes and vice versa
	"os"
)


type Wav struct { //endian
	ChunkID string //big
	ChunkSize uint32 //Little
	Format string //big

	Subchunk1ID string //big
	Subchunk1Size uint32 //Little
	AudioFormat uint16 //Little
	NumChannels uint16 //Little
	SampleRate uint32 //Little
	ByteRate uint32 //Little
	BlockAlign uint16 //Little
	BitsPerSample uint16 //Little

	Subchunk2ID string //big
	Subchunk2Size uint32 //Little
	Data []byte //Little
}

func (wave *Wav) Bytes() []byte { //returns the the Wav as a []byte which can be written to a file
	var bytes []byte
	bytes = append(bytes, []byte(wave.ChunkID)...)
	chunkSizeBytes := []byte{0,0,0,0}
	binary.LittleEndian.PutUint32(chunkSizeBytes, wave.ChunkSize)
	bytes = append(bytes, chunkSizeBytes...)
	bytes = append(bytes, []byte(wave.Format)...)

	bytes = append(bytes, []byte(wave.Subchunk1ID)...)
	subchunk1SizeBytes := []byte{0,0,0,0}
	binary.LittleEndian.PutUint32(subchunk1SizeBytes, wave.Subchunk1Size)
	bytes = append(bytes, subchunk1SizeBytes...)
	audioFormatBytes := []byte{0,0}
	binary.LittleEndian.PutUint16(audioFormatBytes, wave.AudioFormat)
	bytes = append(bytes, audioFormatBytes...)
	numChannelsBytes := []byte{0,0}
	binary.LittleEndian.PutUint16(numChannelsBytes, wave.NumChannels)
	bytes = append(bytes, numChannelsBytes...)
	sampleRateBytes := []byte{0,0,0,0}
	binary.LittleEndian.PutUint32(sampleRateBytes, wave.SampleRate)
	bytes = append(bytes, sampleRateBytes...)
	byteRateBytes := []byte{0,0,0,0}
	binary.LittleEndian.PutUint32(byteRateBytes, wave.ByteRate)
	bytes = append(bytes, byteRateBytes...)
	blockAlignBytes := []byte{0,0}
	binary.LittleEndian.PutUint16(blockAlignBytes, wave.BlockAlign)
	bytes = append(bytes, blockAlignBytes...)
	bitsPerSampleBytes := []byte{0,0}
	binary.LittleEndian.PutUint16(bitsPerSampleBytes, wave.BitsPerSample)
	bytes = append(bytes, bitsPerSampleBytes...)

	bytes = append(bytes, []byte(wave.Subchunk2ID)...)
	subchunk2SizeBytes := []byte{0,0,0,0}
	binary.LittleEndian.PutUint32(subchunk2SizeBytes, wave.Subchunk2Size)
	bytes = append(bytes, subchunk2SizeBytes...)
	bytes = append(bytes, wave.Data...)

	return bytes //returns the contents of a complete .wav file
}

func Create(data []byte, numOfChannels uint16, samplerate uint32, bitDepth uint16) *Wav{ //creates a Wav instance, patches it up
	wave := Wav{
		NumChannels: numOfChannels,
		SampleRate: samplerate,
		BitsPerSample: bitDepth,
		Data: data}
	(&wave).Patch()
	return &wave //!!!YOU CAN DO THIS IN GO WITOUT THE POINTER GETTING CLEANED UP (because garbage collection)
}

func (wave *Wav) Patch() { //sets all the variables in the Wav to the correct values according to the Wav's current NumChannels, SampleRate, BitsPerSample and Data
	wave.ChunkID = "RIFF" //THERE IS NO -> OPERATOR JUST . IS FINE APARENTLY
	wave.Format = "WAVE"
	wave.Subchunk1ID = "fmt "
	wave.Subchunk2ID = "data"

	wave.Subchunk1Size = 16
	wave.AudioFormat = 1

	wave.ByteRate = wave.SampleRate * uint32(wave.NumChannels) * uint32(wave.BitsPerSample/8)
	wave.BlockAlign = wave.NumChannels * (wave.BitsPerSample/8)
	wave.Subchunk2Size = uint32(len(wave.Data))
	wave.ChunkSize = 36 + wave.Subchunk2Size
}

func (wave *Wav) Save(fileLocation string) error { //saves a Wav to fileLocation
	f, err := os.Create(fileLocation)
	bytes := wave.Bytes()
	f.Write(bytes)
	f.Close()
	return err
}

func Read(fileLocation string) (*Wav, error) { //creates a Wav from the file at fileLocation
	buffer := make([]byte, 8)
	f, err := os.Open(fileLocation)
	f.Read(buffer)
	size := binary.LittleEndian.Uint32(buffer[4:]) //2nd 4th-8th bytes is the chunkSize, convert it to an int
	dataBuffer := make([]byte, size)
	f.Read(dataBuffer) //dataBuffer IS EVERYTHING AFTER THE FIRST 8 BYTES
	f.Close()
	channels := binary.LittleEndian.Uint16(dataBuffer[14:16]) //REMEMBER, THE FIRST 8 BYTES OF THE FILE ARE GONE
	samplerate := binary.LittleEndian.Uint32(dataBuffer[16:20])
	bitDepth := binary.LittleEndian.Uint16(dataBuffer[26:28])
	return Create(dataBuffer[36:], channels, samplerate, bitDepth), err //first 36 bytes is still meta data!
}


func main(){
	wave, _ := Read("C:\\Users\\Wakydawgster\\Desktop\\DigitWhore.wav")
	fmt.Println(wave.NumChannels, "channels, sample rate of", wave.SampleRate, "hz,", wave.BitsPerSample, "bit depth")
	fmt.Println("Writing", wave.ChunkSize + 8, "bytes...")
	wave.Save("C:\\Users\\Wakydawgster\\Desktop\\meme69_tro11er.wav")
	//wave.Data = []byte("Ze big meme")
	//wave.Patch()
	//wave.Save("C:\\Users\\Wakydawgster\\Desktop\\meme69_tro11er.wav")
	fmt.Println("Complete")
}


/*
func Create(data []byte, numOfChannels uint16, samplerate uint32, bitDepth uint16) Wav{ //creates a Wav instance
	wave := Wav{
		ChunkID: "RIFF", 
		Format: "WAVE", 
		Subchunk1ID: "fmt ", 
		Subchunk2ID: "data",

		Subchunk1Size: 16,
		AudioFormat: 1,

		NumChannels: numOfChannels,
		SampleRate: samplerate,
		BitsPerSample: bitDepth,
		ByteRate: samplerate * uint32(numOfChannels) * uint32(bitDepth/8),
		BlockAlign: numOfChannels * (bitDepth/8),
	
		Subchunk2Size: uint32(len(data)),
		ChunkSize: 36 + uint32(len(data)),
		Data: data}
	return wave
}
*/