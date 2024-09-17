function handleLikeDislike(postId, action) {
    const likeButton = document.getElementById(`like-button-${postId}`);
    const dislikeButton = document.getElementById(`dislike-button-${postId}`);
    const likeCount = document.getElementById(`like-count-${postId}`);
    const dislikeCount = document.getElementById(`dislike-count-${postId}`);
    const likeLabel = document.getElementById(`like-label-${postId}`);
    const dislikeLabel = document.getElementById(`dislike-label-${postId}`);
    // Disable both buttons temporarily
    likeButton.disabled = true;
    dislikeButton.disabled = true;
    
    // Make the request to the server
    fetch(`/like_dislike`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ postId: parseInt(postId), action: action }),
    })
    .then((response) => response.json())
    .then((data) => {
      // Re-enable buttons after response
      likeButton.disabled = false;
      dislikeButton.disabled = false;
      
      // Update counts based on the response
      likeCount.textContent = data.likeCount;
      dislikeCount.textContent = data.dislikeCount;
  
      // Reset styles for both buttons first
      likeButton.classList.remove('text-blue-500');
      dislikeButton.classList.remove('text-red-500');
      likeLabel.classList.remove('text-blue-500');
      dislikeLabel.classList.remove('text-red-500');
      
      // Apply new styles based on the reaction
      if (data.action === 'like') {
        likeButton.classList.add('text-blue-500'); // Highlight like button
        likeLabel.classList.add('text-blue-500'); // Highlight like label
      } else if (data.action === 'dislike') {
        dislikeButton.classList.add('text-red-500'); // Highlight dislike button
        dislikeLabel.classList.add('text-red-500'); // Highlight dislike label
      }else{
        ikeButton.classList.remove('text-blue-500');
      dislikeButton.classList.remove('text-red-500');
      likeLabel.classList.remove('text-blue-500');
      dislikeLabel.classList.remove('text-red-500');
      }
    })
    .catch((error) => {
      console.error('Error:', error);
      // Re-enable buttons if an error occurs
      likeButton.disabled = false;
      dislikeButton.disabled = false;
    });
}