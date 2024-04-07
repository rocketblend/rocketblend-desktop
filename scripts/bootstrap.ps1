docker build -t rocketblend_desktop_builder .
docker run --rm -v ${PWD}/build/bin:/artifacts rocketblend_desktop_builder