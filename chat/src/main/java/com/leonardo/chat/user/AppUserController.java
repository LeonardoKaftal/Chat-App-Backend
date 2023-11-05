package com.leonardo.chat.user;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("api/v1/users")
@RequiredArgsConstructor
public class AppUserController {

    private final AppUserService appUserService;


    // very simple method that is used to check if the client have a valid jwt token on the client side
    @GetMapping("/authenticate")
    public ResponseEntity<String> checkIfAuthenticated() {
        return ResponseEntity.ok().body("You are successfully authenticated");
    }

    @GetMapping("/email/{email}")
    public ResponseEntity<String> getUsernameByUsingEmail(@PathVariable String email) {
        String username = appUserService.loadUserByUsername(email).getUsername();
        return ResponseEntity.ok().body(username);
    }
}
