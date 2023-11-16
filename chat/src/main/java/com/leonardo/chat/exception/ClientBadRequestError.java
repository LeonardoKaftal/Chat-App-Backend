package com.leonardo.chat.exception;

public class ClientBadRequestError extends RuntimeException {
    public ClientBadRequestError(String message) {
        super(message);
    }
}
