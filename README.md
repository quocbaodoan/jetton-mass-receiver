<p align="center" style="color: #444">
  <h1 align="center">🍭 Jetton Mass Receiver</h1>
</p>
<p align="center" style="font-size: 1.2rem;">Send TON Jettons Easily!</p>

1. Prepare json file with the following structure with any file name:

```json
[
  {
    "seed": "example seed phrase 1"
  },
  {
    "seed": "example seed phrase 2"
  }
]
```

2. Clone Repository

```bash
git clone git@git@github.com:quocbaodoan/jetton-mass-receiver.git; cd jetton-mass-receiver;
```

3. [Install golang](https://go.dev/doc/install)
4. Setup Jetton Receiver

```bash
go run src/main.go src/cli.go src/highload.go setup
```

5. Run Jetton Receiver

```bash
go run src/main.go src/cli.go src/highload.go
```
