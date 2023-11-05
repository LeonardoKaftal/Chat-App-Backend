package com.leonardo.chat.exception.handler;

import java.time.LocalDateTime;

public record ApiError(String path, String message, int statusCode, LocalDateTime localDateTime) {
}
