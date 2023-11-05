package com.leonardo.chat.security.auth;

import com.leonardo.chat.security.jwt.TokenService;
import com.leonardo.chat.user.AppUserService;
import lombok.RequiredArgsConstructor;
import org.springframework.lang.NonNull;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.Authentication;
import org.springframework.stereotype.Service;

@Service
@RequiredArgsConstructor
public class AuthService {

    private final TokenService tokenService;
    private final AuthenticationManager authenticationManager;
    private final AppUserService appUserService;
    public String register(@NonNull RegisterRequest registerRequest) {
        appUserService.registerUser(registerRequest);
        Authentication authentication = new UsernamePasswordAuthenticationToken(registerRequest.email(),registerRequest.password());
        return tokenService.generateToken(authentication);
    }

    public String authenticate(@NonNull LoginRequest userLogin) {
        Authentication authentication = authenticationManager.authenticate(new UsernamePasswordAuthenticationToken(userLogin.email(), userLogin.password()));
        return tokenService.generateToken(authentication);
    }
}
