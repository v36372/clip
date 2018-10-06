#!/bin/bash
filename=$1
mf=$2
sf=$3
mt=$4
st=$5
outfile=$6
echo ""
echo "----------------------------------------------"

############ STEP 1
echo "+ STEP 1: BEGIN CUTTING MKV FROM INPUT FILE"
echo ""
ffmpeg -i ./r/${filename}.mkv -c:v copy -c:a aac -ss 00:${mf}:${sf} -t 00:${mt}:${st} -async 1 -y -strict -2 ./vids/${outfile}.mkv
echo ""

########### STEP 2
echo "+ STEP 2: GENERATE THUMBNAIL"
ffmpeg -ss 5 -i ./vids/${outfile}.mkv -frames:v 1 -q:v 2  -filter:v scale="260:150" ./../stream/public/images/${outfile}.jpg
echo ""
#mv new.mp4 old.mp4
#mv new.mkv latest.mkv

echo "FINISHED CUTTING MKV"
echo "-----------"
echo ""

