package com.leonardo.chat.security.jwt;

import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.NoSuchAlgorithmException;

public class KeyGeneratorUtils {
    private KeyGeneratorUtils() {}

    public static KeyPair generateRsaKey() throws NoSuchAlgorithmException {
        KeyPair keyPair;

            KeyPairGenerator keyPairGenerator = KeyPairGenerator.getInstance("RSA");
            keyPairGenerator.initialize(2048);
            keyPair = keyPairGenerator.generateKeyPair();

        return keyPair;
    }
}
