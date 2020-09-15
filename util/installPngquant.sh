# shellcheck disable=SC1113
#/bin/bash
png='pngquant'
type $png >/dev/null 2>&1 || {
  echo >&2 "$png not installed. now install";
  osname=`uname -s`
  echo "os:$osname"
  case $osname in
  "Darwin" | "darwin")
    echo "darwin"
    str=$(brew install $png)
    echo ${#str}
    ;;
  "Linux" | "linux")
    echo "linux"
    str=$(yum install $png)
    echo ${#str}
    ;;
  *)
    echo "unknown $osname";;
  esac
  exit 1;
}