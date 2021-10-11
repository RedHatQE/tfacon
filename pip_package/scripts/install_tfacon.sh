#!/bin/bash
mkdir -p ~/.local/bin/
cp tfacon_pip/tfacon_binary/tfacon ~/.local/bin/
export PATH="$HOME/.local/bin/:$PATH"