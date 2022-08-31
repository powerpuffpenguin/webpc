Target="webpc"
Docker="king011/webpc"
Dir=$(cd "$(dirname $BASH_SOURCE)/.." && pwd)
Version="v1.2.8"
View=1
Platforms=(
    windows/amd64
    darwin/amd64
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
    forward/logger/logger.proto
    forward/system/system.proto
    forward/shell/shell.proto
    forward/vnc/vnc.proto
    forward/forward/forward.proto
)
