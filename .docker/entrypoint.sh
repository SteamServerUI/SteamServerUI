#!/usr/bin/env sh

cp /opt/SSUIBuildFiles/SteamServerUI /app/SteamServerUI
cp /opt/SSUIBuildFiles/LICENSE /app/LICENSE
chmod +x /app/SteamServerUI
exec /app/SteamServerUI "$@"