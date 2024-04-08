CAMERA=/dev/video0
FILE=`date +%s`

# Basic auth user / password
TOKEN=`printf "user:pass" | base64`
URL="https://server.com/send"

mkdir /home/pi/pics
cd /home/pi/pics

# take picture
ffmpeg -hide_banner -loglevel error -y -f video4linux2 -s 1280x1024 -i ${CAMERA} -frames 1 pic.jpeg

# add timestamp on image
ffmpeg -hide_banner -loglevel error -y -i pic.jpeg -vf "drawtext=text='%{localtime}': x=(w-tw)/2: y=h-(2*lh): fontcolor=white: box=1: boxcolor=0x00000000@1: fontsize=30" -r 25 -t 5 ${FILE}.jpeg

# upload
curl -H"Authorization: Basic ${TOKEN}" -F "file=@./${FILE}.jpeg" ${URL}
