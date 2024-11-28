const ffmpeg = require('fluent-ffmpeg');
const fs = require('fs');
const { fft } = require('fft-js');

const filePath = '../audio_files/audio_1.mp3'; // Replace with your file path
const tempFilePath = 'output_audio_data.raw';

// Convert audio to raw PCM data
function convertToPCM(inputFile, outputFile) {
    return new Promise((resolve, reject) => {
      ffmpeg(inputFile)
        .output(outputFile)
        .audioChannels(1) // Mono channel
        .audioFrequency(44100) // Sample rate
        .format('s16le') // PCM 16-bit little-endian
        .on('end', () => {
          // Check if the output file is valid
          if (!fs.existsSync(outputFile) || fs.statSync(outputFile).size === 0) {
            return reject(new Error('Output PCM file is empty or invalid.'));
          }
          resolve();
        })
        .on('error', (err) => reject(err))
        .run();
    });
  }

// Generate FFT from raw PCM data
function generateFFT(filePath) {
    console.log('Generating FFT from raw PCM data...', filePath);
    const buffer = fs.readFileSync(filePath);
    const samples = new Int16Array(buffer); // Convert buffer to 16-bit integer array

    // Check if samples are valid
    if (!samples || samples.length === 0) {
      throw new Error('Audio data is empty or invalid.');
    }

    // Normalize data to a range between -1 and 1 (optional but recommended)
    const normalizedSamples = Array.from(samples).map((val) => val / 32768);

    // Compute FFT
    const fftData = fft(normalizedSamples);
    return fftData.magnitude; // Get magnitude values of FFT
  }

(async () => {
  try {
    console.log('Converting audio to PCM...');
    await convertToPCM(filePath, tempFilePath);

    console.log('Generating FFT...');
    const fftValues = generateFFT(tempFilePath);

    // Save FFT values to JSON file
    const fftJson = {
      file: filePath,
      fft: fftValues,
    };

    fs.writeFileSync('fft_output.json', JSON.stringify(fftJson, null, 2));
    console.log('FFT data saved to fft_output.json');

    // Clean up temporary file
    fs.unlinkSync(tempFilePath);
  } catch (error) {
    console.error('Error:', error);
  }
})();
