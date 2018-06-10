package main

import (
	"fmt"

	"github.com/robertkrimen/otto"
)

func main() {
	vm := otto.New()

	result, _ := vm.Run(`
		getHoney {
			var t = Math.floor((new Date).getTime() / 1e3),
				e = t.toString(16).toUpperCase(),
				i = CryptoJS.MD5(t + '').toString().toUpperCase();
			if (8 != e.length) return {
				as: "479BB4B7254C150",
				cp: "7E0AC8874BB0985"
			};
			for (var n = i.slice(0, 5), a = i.slice(-5), s = "", o = 0; 5 > o; o++)
				s += n.substr(o, 1) + e.substr(o, 1);
			for (var r = "", c = 0; 5 > c; c++)
				r += e.substr(c + 3, 1) + a.substr(c, 1);
			return {
				as: "A1" + s + e.slice(-3),
				cp: e.slice(0, 3) + r + "E1"
			}
		}
	`)
	fmt.Println(result)

}
