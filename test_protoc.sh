pushd /home/supertapood/Desktop/yapbc
cargo run -- compile -w ../Flint/protobuf -f . -o ../flint/core/generated -l go -m github.com/SuperTapood/Flint/core/generated
cargo run -- compile -w ../Flint/protobuf -f . -o ../pyflint/generated -l python
popd