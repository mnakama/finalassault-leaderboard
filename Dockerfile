FROM alpine
#FROM linuxkit/ca-certificates:c1c73ef590dffb6a0138cf758fe4a4305c9864f4

RUN /sbin/apk add --no-cache tzdata ca-certificates mailcap

ADD leaderboard.html /
ADD res /res
ADD style /style
ADD finalassault-leaderboard /

USER 900

ENTRYPOINT ["/finalassault-leaderboard"]