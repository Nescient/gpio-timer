# gpio-timer
A GPIO timer for DebyNet server (https://derbynet.org/)

## Libre Computer Wiring
https://hub.libre.computer/t/libre-computer-wiring-tool/40

## Provisioning
```
sudo apt update

wget https://github.com/Kitware/CMake/releases/download/v3.26.4/cmake-3.26.4-linux-x86_64.sh
sh cmake-3.26.4-linux-x86_64.sh

# sudo apt install cmake # why is the package manager so old?!?!!?

sudo apt install build-essential clang clang-format clang-tidy cppcheck doxygen codespell lcov python3-pip curl libcurl4-openssl-dev
pip install conan cmake-init
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
```

## Third-Party Libs
https://github.com/friendlyanon/cmake-init

## Building and installing

See the [BUILDING](BUILDING.md) document.

## Contributing

See the [CONTRIBUTING](CONTRIBUTING.md) document.

## Licensing

<!--
Please go to https://choosealicense.com/licenses/ and choose a license that
fits your needs. The recommended license for a project of this type is the
GNU AGPLv3.
-->