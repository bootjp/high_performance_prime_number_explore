#include <stdio.h>
#include <math.h>

int main() {
    unsigned n;
    int count = 2;
    scanf("%d", &n);
    printf("1\t : 2\n");
    for (unsigned i = 3; i <= n; i += 2) {
        int isPrime = 1;
        if (i % 2 == 0) {
            isPrime = 0;
        } else {
            for (unsigned j = 3; j*j <= i; j += 2) {
                if (i % j == 0) {
                    isPrime = 0;
                    break;
                }
            }
        }
        if (isPrime == 1) {
            printf("%d\t : %d\n", count, i);
            count++;
        }
    }
    return 0;
}
