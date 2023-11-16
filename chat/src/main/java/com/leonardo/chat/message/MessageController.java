package com.leonardo.chat.message;

import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.lang.NonNull;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/v1/messages")
@RequiredArgsConstructor
public class MessageController {

    private final MessageService messageService;

    @GetMapping("/{roomName}")
    public ResponseEntity<List<MessageResponse>> getAllMessagesForARoom(@NonNull @PathVariable String roomName) {
        return ResponseEntity.ok(messageService.getAllMessageInARoom(roomName));

    }

    @PostMapping
    public ResponseEntity<List<MessageResponse>> saveMessage(@NonNull @RequestBody MessageRequest messageRequest) {
        return ResponseEntity.ok(messageService.saveMessage(messageRequest));
    }

}
