 Use these styles  for the code after it in tailwind cdn: 

 ```
### **Colors**
1. **Background Colors**:
   - `bg-[#1a1a1a]`: Very dark gray for containers, forms, and cards.
   - `bg-[#2a2a2a]`: Slightly lighter dark gray for input fields and interactive elements.

2. **Text Colors**:
   - `text-[#e5e5e5]`: Light gray for primary text.
   - `text-[#a1a1aa]`: Lighter gray for secondary text (e.g., placeholders, metadata).
   - `text-[#4C9C94]`: Greenish-teal for icons and interactive elements.
   - `text-red-500`: Red for delete buttons and error states.

3. **Border Colors**:
   - `border-[#3a3a3c]`: Medium-dark gray for borders.
   - `hover:border-[#4C9C94]`: Greenish-teal for hover states on borders.

4. **Interactive Element Colors**:
   - `text-[#4C9C94]`: Greenish-teal for icons and buttons.
   - `hover:text-[#3a7a6f]`: Darker greenish-teal for hover states.
   - `hover:bg-[#3a7a6f]`: Darker greenish-teal for button hover backgrounds.

---

### **Typography**
1. **Font Sizes**:
   - `text-lg`: Large text for headings.
   - `text-sm`: Small text for subheadings or metadata.

2. **Font Weights**:
   - `font-semibold`: Semi-bold for headings.
   - `font-medium`: Medium for subheadings.

---

### **Spacing**
1. **Padding**:
   - `p-6`: 1.5rem (24px) for containers.
   - `p-4`: 1rem (16px) for cards and panels.
   - `p-3`: 0.75rem (12px) for list items.
   - `p-2`: 0.5rem (8px) for buttons and input fields.

2. **Margin**:
   - `mb-4`: Margin-bottom of 1rem (16px) for spacing between sections.
   - `mt-4`: Margin-top of 1rem (16px) for spacing between sections.

3. **Spacing Between Elements**:
   - `space-x-2`: Horizontal spacing of 0.5rem (8px) between child elements.
   - `space-x-3`: Horizontal spacing of 0.75rem (12px) between child elements.
   - `space-y-2`: Vertical spacing of 0.5rem (8px) between child elements.
   - `space-y-4`: Vertical spacing of 1rem (16px) between child elements.

---

### **Borders**
1. **Border Styles**:
   - `border`: Default border width.
   - `border-b`: Bottom border.
   - `border-2`: Border width of 2px.
   - `border-dashed`: Dashed border style.

2. **Border Radius**:
   - `rounded-lg`: Large border radius (0.5rem or 8px) for rounded corners.

---

### **Transitions and Animations**
1. **Transitions**:
   - `transition-colors duration-200`: Smooth color transitions with a duration of 200ms.
   - `transition-opacity duration-200`: Smooth opacity transitions with a duration of 200ms.
   - `transition-all duration-200`: Smooth transitions for all properties with a duration of 200ms.

2. **Animations**:
   - `animate-bounce`: Bounce animation for interactive elements.
   - `group-hover:animate-pulse`: Pulse animation on hover for elements within a group.
   - `group-active:animate-ping`: Ping animation on active state for elements within a group.

---

### **Flexbox Layout**
1. **Flex Containers**:
   - `flex`: Flexbox container.
   - `flex-col`: Flexbox column layout.
   - `flex-grow`: Allow element to grow and fill available space.

2. **Alignment**:
   - `items-center`: Center items vertically.
   - `justify-center`: Center items horizontally.

---

### **Interactive Elements**
1. **Buttons**:
   - Use `text-[#4C9C94]` for primary buttons and `hover:text-[#3a7a6f]` for hover states.
   - Use `text-red-500` for delete buttons and `hover:text-red-700` for hover states.
   - Add `transition-colors duration-200` for smooth hover effects.

2. **Hover Effects**:
   - `hover:bg-[#3a3a3c]`: Darker background on hover.
   - `hover:border-[#4C9C94]`: Greenish-teal border on hover.
   - `group-hover:opacity-100`: Show elements on hover within a group.

---

### **Scrollbars**
1. **Scrollbar Styling**:
   - `overflow-y-auto`: Enable vertical scrolling.
   - `custom-scrollbar`: Custom scrollbar styling (requires additional CSS).

---

### **Component-Specific Styles**
1. **Containers**:
   - Background: `bg-[#1a1a1a]`
   - Text: `text-[#e5e5e5]`
   - Border: `border-[#3a3a3c]`
   - Padding: `p-6` or `p-4` depending on the context.

2. **Input Fields**:
   - Background: `bg-[#2a2a2a]`
   - Text: `text-[#e5e5e5]`
   - Placeholder: `placeholder-[#a1a1aa]`
   - Focus Border: `focus:border-[#4C9C94]`

3. **Icons**:
   - Use `text-[#4C9C94]` for icons.
   - Add hover effects with `hover:text-[#3a7a6f]`.

---

### **Summary of Primary Colors**
- **Dark Gray**: `#1a1a1a`, `#2a2a2a`, `#3a3a3c`
- **Light Gray**: `#e5e5e5`, `#a1a1aa`
- **Greenish-Teal**: `#4C9C94`, `#3a7a6f`
- **Red**: `#ef4444` (Tailwind’s `red-500`), `#dc2626` (Tailwind’s `red-700`)

---

### **How to Use Consistently**
1. **Backgrounds**: Use `bg-[#1a1a1a]` for containers and `bg-[#2a2a2a]` for input fields.
2. **Text**: Use `text-[#e5e5e5]` for primary text and `text-[#a1a1aa]` for secondary text.
3. **Borders**: Use `border-[#3a3a3c]` for default borders and `hover:border-[#4C9C94]` for hover states.
4. **Buttons**: Use `text-[#4C9C94]` for primary buttons and `text-red-500` for delete buttons.
5. **Transitions**: Always add `transition-colors duration-200` for interactive elements.
6. **Spacing**: Use `p-4` for cards, `p-2` for buttons, and `space-y-4` for vertical spacing.

 <!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Chat Application</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/htmx.org@2.0.4"></script>
    <script src="https://unpkg.com/htmx-ext-sse@2.2.2/sse.js"></script>
</head>

<body class="bg-[#1a1a1a] text-[#e5e5e5] p-6 flex flex-col h-screen">

    <!-- Display streaming response -->
    <div
        id="main-content"
        class="flex-grow p-8 mt-16 overflow-y-auto scrollbar-thin scrollbar-thumb-[#4C9C94] scrollbar-track-[#1a1a1a] border border-[#3a3a3c] rounded-lg mb-4"
    >
        <div 
            id="stream-response" 
            hx-ext="sse" 
            sse-connect="/stream" 
            sse-swap="update" 
            hx-swap="beforeend">
            <!-- Responses will be appended here -->
        </div>
    </div>

    <!-- Prompt area -->
    <div class="p-6 flex flex-col">
        <div class="text-[#e5e5e5] flex-grow flex flex-col">
            <form
                hx-post="/prompt" 
                hx-swap="none" 
                hx-on:htmx:before-request="document.getElementById('stream-response').innerHTML = '';"
                class="flex items-center space-x-2 bg-[#1a1a1a] rounded-lg border border-[#3a3a3c] p-2 hover:border-[#4C9C94] transition-colors duration-200"
            >
                <!-- Send Icon (Interactive Animations) -->
                <button
                    type="submit"
                    class="p-1 text-[#4C9C94] hover:text-[#007acc] transition-colors duration-200 flex items-center justify-center group"
                >
                    <svg
                        class="w-4 h-4 hover:w-5 hover:h-5 transition-all duration-200 animate-bounce group-hover:animate-pulse group-active:animate-ping"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                        xmlns="http://www.w3.org/2000/svg"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M13 5l7 7-7 7M5 5l7 7-7 7"
                        ></path>
                    </svg>
                </button>
                <!-- Input Field -->
                <input
                    type="text"
                    id="prompt-input"
                    name="prompt"
                    placeholder="Type your prompt..."
                    class="w-full p-2 bg-transparent text-[#e5e5e5] focus:outline-none placeholder-[#a1a1aa]"
                    required
                />
                <input type="hidden" id="prompt-index" name="prompt-index" value="-1"/>
            </form>
        </div>
    </div>
</body>

</html>

in the  input are while the tokens are  being generated, make the pormpt appear in redish color and have a stop effect indicating if the use wants to stop the request
