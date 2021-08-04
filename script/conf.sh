Target="webpc"
Docker="github.com/powerpuffpenguin/webpc"
Dir=$(cd "$(dirname $BASH_SOURCE)/.." && pwd)
Version="v0.0.1"
View=1
Platforms=(
    darwin/amd64
    windows/amd64
    linux/arm
    linux/amd64
)
UUID="81dd3b50-f343-11eb-8332-dfc4915441d6"
Protos=(
    system/system.proto
    session/session.proto
    user/user.proto
    logger/logger.proto
    slave/slave.proto
    group/group.proto

    forward/fs/fs.proto
)