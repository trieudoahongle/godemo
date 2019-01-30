#include <iostream>
class Foo {
public:
  int a;
  Foo(int _a){
	  this->a = _a;
  };
  ~Foo(){};
  void increase(int i){
	  a += i;
  };
  int getValue(){
	  return a;
  };
};
int main(){
	Foo f(1);
	return f.getValue();
}
