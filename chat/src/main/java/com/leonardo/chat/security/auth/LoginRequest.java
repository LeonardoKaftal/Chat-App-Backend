package com.leonardo.chat.security.auth;

import lombok.Builder;

@Builder
public record LoginRequest(String email, String password) {
}
