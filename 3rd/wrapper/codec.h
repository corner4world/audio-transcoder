#ifndef CODEC_H
#define CODEC_H

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <assert.h>
#include <string.h>

typedef struct audio_codec
{
    void *handle;
    int sample_rate;
    int sample_size;
    int channels;
    int profile;
    int level;
    char *extra_data;
    int extra_data_size;
} audio_codec_t;

#endif