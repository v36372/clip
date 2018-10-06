#!/bin/bash
filename=$1
S_DIR=r
STREAM_DIR=`date +"%Y%m%d"`
echo ""
echo "----------------------------------------------"

############ STEP 0
echo "MAKE DIRECTORY FOR TODAY STREAM"
mkdir -p $S_DIR
mkdir -p $S_DIR/$STREAM_DIR

############ STEP 1
echo "BEGIN EXTRACTING MKV FROM STREAM CHUNKS"
echo ""
echo "+ STEP 1: GET 10 LATEST STREAM CHUNKS FILES"
a=`ls -Art ./stream/hls | tail -n 1 | cut -d'-' -f 2 | cut -d'.' -f 1`
b=$(($a-10))
if [[ -z "$a"   ]]; then
	echo "   NO STREAM CHUNKS FOUND; ABORT"
	exit 1
fi
if [[ "$b" -lt 0   ]]; then
	echo "   STREAM CHUNKS IS LESS THAN 10 FILES; ABORT"
	exit 1
fi
echo "   STREAM CHUNK FILES COMBINED INTO 1 FILE"
echo ""
for ((i=$b;i<=$a;i++)); do cat ./stream/hls/laptrinhstream-${i}.ts >> $S_DIR/$STREAM_DIR/${a}.ts; done
#ffmpeg -i new.ts -c:v libx264 -c:a copy -bsf:a aac_adtstoasc -y new.mp4
echo ""
echo ""

########### STEP 2
echo "+ STEP 2: ENCODING ts FILES TO MKV $S_DIR/$STREAM_DIR/${a}.mkv"
ffmpeg -i $S_DIR/$STREAM_DIR/${a}.ts -probesize 20000000 -analyzeduration 20000000 -c:v copy -c:a aac -strict -2 -y $S_DIR/$STREAM_DIR/${a}.mkv
echo ""
rm $S_DIR/$STREAM_DIR/${a}.ts
#mv new.mp4 old.mp4
#mv new.mkv latest.mkv

echo "FINISHED EXTRACTING MKV FROM STREAM CHUNKS"
echo "-----------"
echo ""

