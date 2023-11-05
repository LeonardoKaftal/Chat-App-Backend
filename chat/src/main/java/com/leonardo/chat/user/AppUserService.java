package com.leonardo.chat.user;

import com.leonardo.chat.exception.UserAlreadyRegisteredException;
import com.leonardo.chat.security.MyPasswordEncoder;
import com.leonardo.chat.security.auth.RegisterRequest;
import lombok.RequiredArgsConstructor;
import org.springframework.lang.NonNull;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.stereotype.Service;

import static com.leonardo.chat.security.Role.USER;


@Service
@RequiredArgsConstructor
public class AppUserService implements UserDetailsService {
    private final AppUserRepository appUserRepository;
    private final MyPasswordEncoder passwordEncoder;

    public AppUser loadUserByUsingUsername(@NonNull String username) {return appUserRepository.findAppUserByUsername(username)
                .orElseThrow(()-> new UsernameNotFoundException("The user with username " + username + " has not been found"));
    }
    @Override
    public AppUser loadUserByUsername(@NonNull String email) throws UsernameNotFoundException {
        return appUserRepository.findAppUserByEmail(email)
                .orElseThrow(()-> new UsernameNotFoundException("The user with email " + email + " has not been found"));
    }

    public void registerUser(@NonNull RegisterRequest registerRequest) {
        var existingUser = appUserRepository.findAppUserByUsername(registerRequest.username());
        if (existingUser.isPresent()) throw new UserAlreadyRegisteredException("The username " + registerRequest.username() + " is already used");

        existingUser = appUserRepository.findAppUserByEmail(registerRequest.email());
        if (existingUser.isPresent()) throw new UserAlreadyRegisteredException("The email " + registerRequest.email() + " is already used");

        var userToRegister = AppUser.builder()
                .email(registerRequest.email())
                .username(registerRequest.username())
                .password(passwordEncoder.getEncoder().encode(registerRequest.password()))
                .role(USER)
                .build();

        appUserRepository.save(userToRegister);
    }
}
