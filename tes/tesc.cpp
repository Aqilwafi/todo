#include <stdio.h>
#include <stdlib.h>

int w; // Total waktu yang dibutuhkan

void count(int m, int A[], int B[], int T[], int n) {
    for (int i = n; i >= 0; i--) {
        if (m <= 0) break; // Jika kebutuhan sudah terpenuhi

        if (A[i] > 0) {
            if (A[i] >= m) {
                w += T[i] * m; // Hitung waktu
                A[i] -= m; // Kurangi stok
                m = 0; // Kebutuhan terpenuhi
            } else {
                w += T[i] * A[i]; // Gunakan semua stok
                m -= A[i]; // Kurangi kebutuhan
                A[i] = 0; // Stok habis
            }
        }
    }
}

int main() {
    int N; // jumlah level kematangan
    scanf("%d", &N);
    
    if (N < 2) return 0; // N harus minimal 2

    int T[N - 1];  // waktu memasak tiap level
    int A[N]; // stok
    int B[N]; // diperlukan

    // Input waktu memasak
    for (int i = 0; i < N - 1; i++) {
        scanf("%d", &T[i]);
        if (T[i] < 1) return 0; // Validasi input
    }

    // Input stok
    for (int i = 0; i < N; i++) {
        scanf("%d", &A[i]);
        if (A[i] > 1000) return 0; // Validasi input
    }

    // Input kebutuhan
    for (int i = 0; i < N; i++) {
        scanf("%d", &B[i]);
        if (B[i] > 1000) return 0; // Validasi input
    }

    // Cek kebutuhan level 0
    if (A[0] < B[0]) {
        printf("-1\n");
    } else {
        // Proses setiap level
        for (int i = 0; i < N; i++) {
            int K = B[i] - A[i]; // Kebutuhan yang tidak terpenuhi
            if (K <= 0) {
                A[i] -= B[i]; // Kurangi stok
                B[i] = 0; // Kebutuhan sudah terpenuhi
            } else {
                count(K, A, B, T, i - 1); // Panggil count untuk sisa kebutuhan
                A[i] = 0; // Set stok menjadi 0 setelah dihitung
            }   
        }
        printf("%d\n", w); // Output total waktu
    }
    
    return 0;
}

