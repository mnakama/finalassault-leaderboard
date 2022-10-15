FROM alpine

RUN /sbin/apk add --no-cache tzdata ca-certificates mailcap
#FROM linuxkit/ca-certificates:c1c73ef590dffb6a0138cf758fe4a4305c9864f4

ADD leaderboard.html /
ADD res /res
ADD style /style
ADD finalassault-leaderboard /

ENTRYPOINT ["/finalassault-leaderboard"]