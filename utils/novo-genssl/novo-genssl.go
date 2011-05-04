/* vim: set sw=4 sts=4 et foldmethod=syntax : */

/*
 * Copyright (c) 2011 Alexander Færøy <ahf@0x90.dk>
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * * Redistributions of source code must retain the above copyright notice, this
 *   list of conditions and the following disclaimer.
 *
 * * Redistributions in binary form must reproduce the above copyright notice,
 *   this list of conditions and the following disclaimer in the documentation
 *   and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package main

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "flag"
    "fmt"
    "os"
    "time"
)

var hostname = flag.String("hostname", "", "Hostname to generate certificate for")
var bits = flag.Int("bits", 2048, "Bit length for the certificate")
var organization = flag.String("organization", "Goatse Inc.", "Organization for the certificate")

func ErrorF(s string, arguments... interface{}) {
    fmt.Printf(s, arguments...)
    os.Exit(1)
}

func main() {
    flag.Parse()
    timestamp := time.Seconds()

    if len(*hostname) == 0 {
        ErrorF("Error: Missing -hostname parameter.\n")
    }

    fmt.Printf(">>> Generating Private Key: %s\n", *hostname + ".key.pem")
    private_key, error := rsa.GenerateKey(rand.Reader, *bits)

    if error != nil {
        ErrorF("Error: Unable to generate private key: %s\n", error)
    }

    private_key_file, error := os.Open(*hostname + ".key.pem", os.O_WRONLY | os.O_CREAT, 0600)

    if error != nil {
        ErrorF("Error: Unable to open %s for writing: %s\n", *hostname + ".key.pem", error)
    }

    defer private_key_file.Close()

    pem.Encode(private_key_file, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(private_key)})

    template := x509.Certificate{
        SerialNumber: []byte{0},
        Subject: x509.Name{
            CommonName: *hostname,
            Organization: []string{*organization},
        },
        NotBefore: time.SecondsToUTC(timestamp),
        NotAfter: time.SecondsToUTC(365 * 24 * 60 * 60 + timestamp),
        SubjectKeyId: []byte{1, 3, 3, 7},
        KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
    }

    fmt.Printf(">>> Generating Certificate: %s\n", *hostname + ".pem")
    certificate, error := x509.CreateCertificate(rand.Reader, &template, &template, &private_key.PublicKey, private_key)

    if error != nil {
        ErrorF("Error: Unable to generate certificate: %s\n", error)
    }

    certificate_file, error := os.Open(*hostname + ".pem", os.O_WRONLY | os.O_CREAT, 0644)

    if error != nil {
        ErrorF("Error: Unable to open %s for writing: %s\n", *hostname + ".pem", error)
    }

    defer certificate_file.Close()

    pem.Encode(certificate_file, &pem.Block{Type: "CERTIFICATE", Bytes: certificate})
}
