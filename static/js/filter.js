document.addEventListener('DOMContentLoaded', function () {
    // Add event listeners to each category button
    document.querySelectorAll('.category-button').forEach(function (button) {
        button.addEventListener('click', function (event) {
            handleCategoryClick(event, button.textContent.trim());
        });
    });

    // Add event listeners to each post card to make them clickable
    document.querySelectorAll('.card').forEach(function (card) {
        card.addEventListener('click', function () {
            const postId = card.getAttribute('data-post-id');
            window.location.href = `/post-details?post_id=${postId}`;
        });
    });

    function handleCategoryClick(event, category) {
        event.preventDefault(); // Prevent the form's default action

        // Show only posts that match the selected category
        document.querySelectorAll('.card').forEach(function (card) {
            const postCategories = card.getAttribute('data-categories').split(',');
            if (category === 'All' || postCategories.includes(category)) {
                card.style.display = 'block'; // Show the post
            } else {
                card.style.display = 'none'; // Hide the post
            }
        });

        // Update button active states
        document.querySelectorAll('.category-button').forEach(button => {
            button.classList.remove('active');
        });
        event.currentTarget.classList.add('active');
    }
});




// to remove the listeners


document.addEventListener('DOMContentLoaded', function () {
    // Remove existing listeners if any by cloning nodes
    document.querySelectorAll('.like-button').forEach(function (likeLabel) {
        likeLabel.replaceWith(likeLabel.cloneNode(true));
    });

    document.querySelectorAll('.dislike-button').forEach(function (dislikeLabel) {
        dislikeLabel.replaceWith(dislikeLabel.cloneNode(true));
    });

    // Attach event listeners to like and dislike buttons
    document.querySelectorAll('.like-button').forEach(function (likeLabel) {
        likeLabel.addEventListener('click', function (event) {
            event.preventDefault(); // Prevent any default action
            event.stopPropagation(); // Prevents the event from bubbling up
            const postId = likeLabel.id.split('-')[2];
            handleLikeDislike(postId, 'like');
        });
    });

    document.querySelectorAll('.dislike-button').forEach(function (dislikeLabel) {
        dislikeLabel.addEventListener('click', function (event) {
            event.preventDefault(); // Prevent any default action
            event.stopPropagation(); // Prevents the event from bubbling up
            const postId = dislikeLabel.id.split('-')[2];
            handleLikeDislike(postId, 'dislike');
        });
    });

    // Set initial active state based on URL parameter
    const urlParams = new URLSearchParams(window.location.search);
    const activeCategory = urlParams.get('category') || 'All';
    document.querySelectorAll('.category-button').forEach(button => {
        if (button.textContent.trim() === activeCategory) {
            button.classList.add('active');
        }
    });
});