package com.leonardo.chat.message;

import com.leonardo.chat.exception.ClientBadRequestError;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@RequiredArgsConstructor
public class MessageService {

    private final MessageRepository messageRepository;

    public List<MessageResponse> getAllMessageInARoom(String roomName) {
        var messages = messageRepository.getMessagesByDestinationRoom(roomName)
                .orElseThrow(()-> new ClientBadRequestError("The room name you have specified does not exist!"));

        return messages
                .stream()
                .map(message -> new MessageResponse(message.getMessageContent(),message.getDestinationRoom(),message.getSenderName()))
                .toList();

    }

    public List<MessageResponse> saveMessage(MessageRequest messageRequest) {
        var message = Message.builder()
                .senderName(messageRequest.SenderName())
                .destinationRoom(messageRequest.DestinationRoom())
                .messageContent(messageRequest.Payload())
                .build();
        messageRepository.save(message);
        return getAllMessageInARoom(message.getDestinationRoom());
    }
}
