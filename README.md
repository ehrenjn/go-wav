# go-wav
Go library for reading and writing .wav files

Wav - Type that contains data representing a .wav file

Useful Wav member variables:
  .Numchannels - Number of channels
  .SampleRate - Sample rate
  .BitsPerSample - bit depth
  .Data - []byte, the actual data of the .wav
 
Wav methods:
  .Bytes() - returns the the Wav as a []byte which can be written to a file
  .Save(fileLocation string) - saves the Wav file to fileLocation
  .Patch() - sets all the variables in the Wav to the correct values according to the Wav's current NumChannels, SampleRate, BitsPerSample and Data

Functions:
  Create(data []byte, numOfChannels uint16, samplerate uint32, bitDepth uint16) - returns a Wav instance, patches it up
  Read(fileLocation string) - returns a Wav instance created from the file at fileLocation
