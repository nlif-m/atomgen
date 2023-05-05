pkgname=atomgen
pkgver=0.0.0
pkgrel=1
pkgdesc="Atom file generating tool"
arch=('x86_64')
url="https://github.com/nlif-m/$pkgname"
license=('GPL')
makedepends=('go' 'make')
source=("$pkgname-$pkgver.tar.gz"::"https://github.com/nlif-m/${pkgname}/archive/refs/tags/${pkgver}.tar.gz")
sha256sums=('e281b308efea607f3849d97a71316c762e2f0ace8638cd520547513739496fa3')

prepare(){
    cd "$pkgname-$pkgver"
    mkdir -p build/
}

build() {
    cd "$pkgname-$pkgver"
    make
}

check() {
    cd "$pkgname-$pkgver"
    go test 
}

package() {
    cd "$pkgname-$pkgver"
    install -Dm755 $pkgname "$pkgdir"/usr/bin/$pkgname
}
