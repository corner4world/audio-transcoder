#include "codec.h"
#include "opus.h"

// 初始化Opus编码器
static OpusEncoder *opus_encoder_create2(int sample_rate, int channels, int application)
{
    int err;
    OpusEncoder *enc = opus_encoder_create(sample_rate, channels, application, &err);

    if (err != OPUS_OK || enc == NULL)
    {
        printf("Failed to create Opus encoder: %s sample_rate: %d channels: %d application: %d\n", opus_strerror(err), sample_rate, channels, application);
        return NULL;
    }

    opus_encoder_ctl(enc, OPUS_SET_BITRATE(64000));
    //    opus_encoder_ctl(enc, OPUS_SET_SIGNAL(OPUS_SIGNAL_MUSIC));
    //    opus_encoder_ctl(enc, OPUS_SET_COMPLEXITY(10));

    return enc;
}

// 编码PCM数据为Opus
static int opus_encoder_encode(OpusEncoder *enc,
                               uint8_t *pcm,
                               int frame_size,
                               uint8_t *data,
                               int data_size)
{
    int bytes_encoded = opus_encode(enc, (const opus_int16*)pcm, frame_size, data, data_size);
    if (bytes_encoded < 0)
    {
        printf("Opus encoding error: %s frame_size: %d\n", opus_strerror(bytes_encoded), frame_size);
        return 0;
    }

    return bytes_encoded;
}

// 销毁编码器
static void opus_encoder_destroy2(OpusEncoder *enc)
{
    opus_encoder_destroy(enc);
}