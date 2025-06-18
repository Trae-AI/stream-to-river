document.addEventListener('DOMContentLoaded', function() {
    // Cache DOM elements for better performance
    const navItems = document.querySelectorAll('.nav-item[data-href]');
    const headings = document.querySelectorAll('.markdown-content h1, .markdown-content h2, .markdown-content h3, .markdown-content h4, .markdown-content h5, .markdown-content h6');
    const contentBody = document.querySelector('.content-body');
    const sidebar = document.querySelector('.sidebar');
    
    // Performance optimization: Use RAF throttling
    let rafId = null;
    let isScrolling = false;
    
    // ID generation for headings
    headings.forEach((heading, index) => {
        if (!heading.id) {
            const text = heading.textContent.trim();
            let id = text.toLowerCase()
                .replace(/\s+/g, '-')
                .replace(/[^\w\-\u4e00-\u9fa5]/g, '')
                .replace(/^-+|-+$/g, '');
            
            // Ensure unique IDs
            const existingElement = document.getElementById(id);
            if (existingElement) {
                id = `${id}-${index}`;
            }
            
            heading.id = id;
            
            // Add smooth entrance animation
            heading.style.opacity = '0';
            heading.style.transform = 'translateY(20px)';
            
            // Animate in with delay based on position
            setTimeout(() => {
                heading.style.transition = 'opacity 0.8s cubic-bezier(0.19, 1, 0.22, 1), transform 0.8s cubic-bezier(0.19, 1, 0.22, 1)';
                heading.style.opacity = '1';
                heading.style.transform = 'translateY(0)';
            }, index * 150);
        }
    });
    
    // smooth scrolling with easing
    function smoothScrollTo(target, duration = 800) {
        if (!target || !contentBody) return;
        
        const targetPosition = target.offsetTop - contentBody.offsetTop - 20;
        const startPosition = contentBody.scrollTop;
        const distance = targetPosition - startPosition;
        const startTime = performance.now();
        
        // Easing function for smooth animation
        function easeOutExpo(t) {
            return t === 1 ? 1 : 1 - Math.pow(2, -10 * t);
        }
        
        function animation(currentTime) {
            const elapsed = currentTime - startTime;
            const progress = Math.min(elapsed / duration, 1);
            const easedProgress = easeOutExpo(progress);
            
            contentBody.scrollTop = startPosition + (distance * easedProgress);
            
            if (progress < 1) {
                requestAnimationFrame(animation);
            }
        }
        
        requestAnimationFrame(animation);
    }
    
    // navigation click handlers with visual feedback
    navItems.forEach(item => {
        item.addEventListener('click', function(e) {
            e.preventDefault();
            
            // Remove all active states with smooth transition
            navItems.forEach(nav => {
                nav.classList.remove('active');
                nav.style.transform = '';
            });
            
            // Add active state with animation
            this.classList.add('active');
            this.style.transform = 'translateX(8px) scale(1.02)';
            
            // Add ripple effect
            createRippleEffect(this, e);
            
            // Smooth scroll to target
            const href = this.getAttribute('data-href');
            if (href && href.indexOf('#') === 0) {
                const targetId = href.substring(1);
                const target = document.getElementById(targetId);
                if (target) {
                    smoothScrollTo(target);
                }
            }
        });
        
        // hover effects
        item.addEventListener('mouseenter', function() {
            if (!this.classList.contains('active')) {
                this.style.transition = 'all 0.4s cubic-bezier(0.19, 1, 0.22, 1)';
                this.style.transform = 'translateX(6px)';
            }
        });
        
        item.addEventListener('mouseleave', function() {
            if (!this.classList.contains('active')) {
                this.style.transform = 'translateX(0)';
            }
        });
    });
    
    // Create ripple effect for better user feedback
    function createRippleEffect(element, event) {
        const ripple = document.createElement('span');
        const rect = element.getBoundingClientRect();
        const size = Math.max(rect.width, rect.height);
        const x = event.clientX - rect.left - size / 2;
        const y = event.clientY - rect.top - size / 2;
        
        ripple.style.cssText = `
            position: absolute;
            width: ${size}px;
            height: ${size}px;
            left: ${x}px;
            top: ${y}px;
            background: rgba(59, 130, 246, 0.3);
            border-radius: 50%;
            transform: scale(0);
            animation: ripple 0.6s cubic-bezier(0.16, 1, 0.3, 1);
            pointer-events: none;
            z-index: 1;
        `;
        
        // Add ripple styles if not exists
        if (!document.getElementById('ripple-styles')) {
            const style = document.createElement('style');
            style.id = 'ripple-styles';
            style.textContent = `
                @keyframes ripple {
                    to {
                        transform: scale(2);
                        opacity: 0;
                    }
                }
                .nav-item {
                    position: relative;
                    overflow: hidden;
                }
            `;
            document.head.appendChild(style);
        }
        
        element.style.position = 'relative';
        element.appendChild(ripple);
        
        // Clean up after animation
        setTimeout(() => {
            if (ripple.parentNode) {
                ripple.parentNode.removeChild(ripple);
            }
        }, 600);
    }
    
    // Optimized scroll handler with intersection observer
    function updateActiveNav() {
        if (!contentBody) return;
        
        const scrollTop = contentBody.scrollTop;
        const viewportHeight = contentBody.clientHeight;
        let currentHeading = null;
        let minDistance = Infinity;
        
        // Find the heading closest to the top of viewport
        headings.forEach(heading => {
            const headingTop = heading.offsetTop - contentBody.offsetTop;
            const distance = Math.abs(scrollTop + 100 - headingTop);
            
            if (headingTop <= scrollTop + viewportHeight * 0.3 && distance < minDistance) {
                minDistance = distance;
                currentHeading = heading;
            }
        });
        
        // Update navigation with smooth transitions
        navItems.forEach(item => {
            const isActive = item.classList.contains('active');
            const shouldBeActive = currentHeading && 
                item.getAttribute('data-href') === `#${currentHeading.id}`;
            
            if (shouldBeActive && !isActive) {
                item.classList.add('active');
                item.style.transform = 'translateX(8px)';
            } else if (!shouldBeActive && isActive) {
            item.classList.remove('active');
                item.style.transform = '';
            }
        });
    }
    
    // Throttled scroll listener with RAF
    function handleScroll() {
        if (!isScrolling) {
            isScrolling = true;
            if (rafId) {
                cancelAnimationFrame(rafId);
            }
            rafId = requestAnimationFrame(() => {
                updateActiveNav();
                isScrolling = false;
            });
        }
    }
    
    // scroll listener with performance optimizations
    if (contentBody) {
        // Use passive listener for better performance
        contentBody.addEventListener('scroll', handleScroll, { passive: true });
        
        // Add momentum scrolling for iOS
        contentBody.style.webkitOverflowScrolling = 'touch';
    }
    
    // Initialize active navigation
    setTimeout(() => {
                    updateActiveNav();
        
        // Handle URL hash on load with animation
        if (window.location.hash) {
            setTimeout(() => {
                const target = document.querySelector(window.location.hash);
                if (target && contentBody) {
                    // Add highlight effect to target
                    target.style.background = 'linear-gradient(135deg, rgba(59, 130, 246, 0.1), rgba(59, 130, 246, 0.05))';
                    target.style.borderRadius = '8px';
                    target.style.padding = '16px';
                    target.style.marginLeft = '-16px';
                    target.style.marginRight = '-16px';
                    target.style.transition = 'all 0.6s cubic-bezier(0.16, 1, 0.3, 1)';
                    
                    smoothScrollTo(target, 1000);
                    
                    // Remove highlight after animation
                    setTimeout(() => {
                        target.style.background = '';
                        target.style.borderRadius = '';
                        target.style.padding = '';
                        target.style.marginLeft = '';
                        target.style.marginRight = '';
                    }, 2000);
                }
            }, 500);
        }
    }, 200);
    
    // mobile navigation toggle with overlay
    function createMobileToggle() {
        if (window.innerWidth <= 768) {
            let mobileToggle = document.getElementById('mobile-nav-toggle');
            let mobileOverlay = document.getElementById('mobile-overlay');
            
            // Create mobile toggle button
            if (!mobileToggle) {
                mobileToggle = document.createElement('button');
                mobileToggle.id = 'mobile-nav-toggle';
                mobileToggle.innerHTML = `
                    <div class="hamburger-lines">
                        <span class="line line1"></span>
                        <span class="line line2"></span>
                        <span class="line line3"></span>
                    </div>
                `;
                mobileToggle.setAttribute('aria-label', '打开导航菜单');
                mobileToggle.setAttribute('aria-expanded', 'false');
                mobileToggle.style.cssText = `
                    position: fixed;
                    top: 16px;
                    right: 16px;
                    z-index: 2001;
                    width: 48px;
                    height: 48px;
                    background: rgba(255, 255, 255, 0.95);
                    border: 1px solid rgba(59, 130, 246, 0.2);
                    border-radius: 12px;
                    cursor: pointer;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    backdrop-filter: blur(20px);
                    -webkit-backdrop-filter: blur(20px);
                    transition: all 0.4s cubic-bezier(0.19, 1, 0.22, 1);
                    box-shadow: 
                        0 4px 20px rgba(0, 0, 0, 0.08),
                        0 1px 3px rgba(0, 0, 0, 0.1),
                        inset 0 1px 0 rgba(255, 255, 255, 0.6);
                    touch-action: manipulation;
                    user-select: none;
                    overflow: hidden;
                `;
                
                // Add hamburger lines styles
                const hamburgerStyles = document.createElement('style');
                hamburgerStyles.textContent = `
                    .hamburger-lines {
                        width: 20px;
                        height: 20px;
                        position: relative;
                        display: flex;
                        flex-direction: column;
                        justify-content: space-between;
                    }
                    
                    .hamburger-lines .line {
                        width: 100%;
                        height: 2px;
                        background: #3b82f6;
                        border-radius: 2px;
                        transition: all 0.4s cubic-bezier(0.19, 1, 0.22, 1);
                        transform-origin: center;
                    }
                    
                    .hamburger-lines .line1 {
                        transform: translateY(0) rotate(0deg);
                    }
                    
                    .hamburger-lines .line2 {
                        opacity: 1;
                        transform: scaleX(1);
                    }
                    
                    .hamburger-lines .line3 {
                        transform: translateY(0) rotate(0deg);
                    }
                    
                    /* Active state */
                    #mobile-nav-toggle.active .hamburger-lines .line1 {
                        transform: translateY(9px) rotate(45deg);
                        background: #ef4444;
                    }
                    
                    #mobile-nav-toggle.active .hamburger-lines .line2 {
                        opacity: 0;
                        transform: scaleX(0);
                    }
                    
                    #mobile-nav-toggle.active .hamburger-lines .line3 {
                        transform: translateY(-9px) rotate(-45deg);
                        background: #ef4444;
                    }
                    
                    #mobile-nav-toggle.active {
                        background: rgba(254, 242, 242, 0.95) !important;
                        border-color: rgba(239, 68, 68, 0.3) !important;
                        box-shadow: 
                            0 8px 30px rgba(239, 68, 68, 0.15),
                            0 2px 6px rgba(239, 68, 68, 0.1),
                            inset 0 1px 0 rgba(255, 255, 255, 0.8) !important;
                        transform: scale(1.05);
                    }
                    
                    #mobile-nav-toggle:hover {
                        transform: translateY(-1px) scale(1.02);
                        box-shadow: 
                            0 8px 25px rgba(59, 130, 246, 0.15),
                            0 3px 8px rgba(59, 130, 246, 0.1),
                            inset 0 1px 0 rgba(255, 255, 255, 0.8);
                        background: rgba(255, 255, 255, 1);
                        border-color: rgba(59, 130, 246, 0.3);
                    }
                    
                    #mobile-nav-toggle.active:hover {
                        background: rgba(254, 242, 242, 1) !important;
                        border-color: rgba(239, 68, 68, 0.4) !important;
                    }
                    
                    /* Ripple effect */
                    #mobile-nav-toggle::before {
                        content: '';
                        position: absolute;
                        top: 50%;
                        left: 50%;
                        width: 0;
                        height: 0;
                        background: radial-gradient(circle, rgba(59, 130, 246, 0.3) 0%, transparent 70%);
                        border-radius: 50%;
                        transform: translate(-50%, -50%);
                        transition: all 0.6s cubic-bezier(0.19, 1, 0.22, 1);
                        opacity: 0;
                        pointer-events: none;
                    }
                    
                    #mobile-nav-toggle:active::before {
                        width: 60px;
                        height: 60px;
                        opacity: 1;
                        transition: all 0.2s cubic-bezier(0.19, 1, 0.22, 1);
                    }
                `;
                document.head.appendChild(hamburgerStyles);
                
                document.body.appendChild(mobileToggle);
            }
            
            // Create mobile overlay
            if (!mobileOverlay) {
                mobileOverlay = document.createElement('div');
                mobileOverlay.id = 'mobile-overlay';
                mobileOverlay.className = 'mobile-overlay';
                document.body.appendChild(mobileOverlay);
            }
            
            // toggle functionality
            function toggleMobileNav() {
                const isOpen = sidebar.classList.contains('mobile-open');
                
                if (isOpen) {
                    closeMobileNav();
                } else {
                    openMobileNav();
                }
            }
            
            function openMobileNav() {
                sidebar.classList.add('mobile-open');
                mobileOverlay.classList.add('active');
                mobileToggle.classList.add('active');
                
                // Prevent body scroll on mobile
                document.body.style.overflow = 'hidden';
                
                // Add smooth animation
                setTimeout(() => {
                    mobileToggle.style.transform = 'rotate(90deg) scale(1.1)';
                }, 100);
            }
            
            function closeMobileNav() {
                sidebar.classList.remove('mobile-open');
                mobileOverlay.classList.remove('active');
                mobileToggle.classList.remove('active');
                mobileToggle.style.transform = 'rotate(0deg) scale(1)';
                
                // Restore body scroll
                document.body.style.overflow = '';
            }
            
            // Event listeners
            mobileToggle.addEventListener('click', toggleMobileNav);
            
            // Close when clicking overlay
            mobileOverlay.addEventListener('click', closeMobileNav);
            
            // Close when clicking sidebar links
            const sidebarLinks = sidebar.querySelectorAll('.nav-item[data-href]');
            sidebarLinks.forEach(link => {
                link.addEventListener('click', () => {
                    setTimeout(closeMobileNav, 300); // Delay to see the click effect
                });
            });
            
            // Touch gesture support
            let touchStartX = 0;
            let touchStartY = 0;
            let touchEndX = 0;
            let touchEndY = 0;
            
            // Swipe to open sidebar
            document.addEventListener('touchstart', (e) => {
                touchStartX = e.changedTouches[0].screenX;
                touchStartY = e.changedTouches[0].screenY;
            }, { passive: true });
            
            document.addEventListener('touchend', (e) => {
                touchEndX = e.changedTouches[0].screenX;
                touchEndY = e.changedTouches[0].screenY;
                handleSwipeGesture();
            }, { passive: true });
            
            function handleSwipeGesture() {
                const deltaX = touchEndX - touchStartX;
                const deltaY = Math.abs(touchEndY - touchStartY);
                const minSwipeDistance = 80;
                const maxVerticalDistance = 100;
                
                // Horizontal swipe with minimal vertical movement
                if (Math.abs(deltaX) > minSwipeDistance && deltaY < maxVerticalDistance) {
                    const isOpen = sidebar.classList.contains('mobile-open');
                    
                    // Swipe right to open (only if starting from left edge)
                    if (deltaX > 0 && touchStartX < 50 && !isOpen) {
                        openMobileNav();
                    }
                    // Swipe left to close
                    else if (deltaX < 0 && isOpen) {
                        closeMobileNav();
                    }
                }
            }
            
        } else {
            // Remove mobile elements on larger screens
            const mobileToggle = document.getElementById('mobile-nav-toggle');
            const mobileOverlay = document.getElementById('mobile-overlay');
            
            if (mobileToggle) {
                mobileToggle.remove();
            }
            if (mobileOverlay) {
                mobileOverlay.remove();
            }
            
            // Ensure sidebar is visible and body scroll is restored
            sidebar.classList.remove('mobile-open');
            document.body.style.overflow = '';
        }
    }
    
    // Initialize mobile toggle
    createMobileToggle();
    
    // Handle window resize
    window.addEventListener('resize', () => {
        createMobileToggle();
        
        // Reset transforms on resize
        if (window.innerWidth > 768) {
            sidebar.classList.remove('mobile-open');
            const mobileToggle = document.getElementById('mobile-nav-toggle');
            if (mobileToggle) {
                mobileToggle.remove();
            }
        }
    });
    
    // keyboard navigation
    document.addEventListener('keydown', (e) => {
        // ESC to close mobile menu
        if (e.key === 'Escape' && sidebar.classList.contains('mobile-open')) {
            sidebar.classList.remove('mobile-open');
            const mobileToggle = document.getElementById('mobile-nav-toggle');
            if (mobileToggle) {
                mobileToggle.classList.remove('active');
                mobileToggle.style.transform = 'rotate(0deg) scale(1)';
            }
        }
        
        // Arrow keys for navigation
        if (e.key === 'ArrowUp' || e.key === 'ArrowDown') {
            const activeNav = document.querySelector('.nav-item.active');
            if (activeNav) {
                const allNavItems = Array.from(navItems);
                const currentIndex = allNavItems.indexOf(activeNav);
                let nextIndex;
                
                if (e.key === 'ArrowUp') {
                    nextIndex = currentIndex > 0 ? currentIndex - 1 : allNavItems.length - 1;
                } else {
                    nextIndex = currentIndex < allNavItems.length - 1 ? currentIndex + 1 : 0;
                }
                
                if (allNavItems[nextIndex]) {
                    allNavItems[nextIndex].click();
                    allNavItems[nextIndex].focus();
                }
                
                e.preventDefault();
            }
        }
    });
    
    // Add loading animation for images
    const images = document.querySelectorAll('img');
    images.forEach(img => {
        if (!img.complete) {
            img.style.opacity = '0';
            img.style.transform = 'scale(0.95)';
            img.style.transition = 'opacity 0.6s cubic-bezier(0.16, 1, 0.3, 1), transform 0.6s cubic-bezier(0.16, 1, 0.3, 1)';
            
            img.addEventListener('load', () => {
                img.style.opacity = '1';
                img.style.transform = 'scale(1)';
            });
        }
    });
    
    // Wrap tables in responsive containers
    const tables = document.querySelectorAll('.markdown-content table');
    tables.forEach(table => {
        // Check if table is not already wrapped
        if (!table.parentNode.classList.contains('table-wrapper')) {
            const wrapper = document.createElement('div');
            wrapper.className = 'table-wrapper';
            table.parentNode.insertBefore(wrapper, table);
            wrapper.appendChild(table);
        }
    });
    
    // Performance monitoring
    if (window.performance && window.performance.mark) {
        window.performance.mark('navigation-enhanced-complete');
    }
    
    // Add subtle parallax effect to content
    if (contentBody && window.innerWidth > 768) {
        let parallaxRaf = null;
        
        contentBody.addEventListener('scroll', () => {
            if (parallaxRaf) {
                cancelAnimationFrame(parallaxRaf);
            }
            
            parallaxRaf = requestAnimationFrame(() => {
                const scrolled = contentBody.scrollTop;
                const rate = scrolled * -0.02;
                
                // Apply subtle parallax to headings
                headings.forEach((heading, index) => {
                    if (index === 0) { // Only apply to main heading
                        heading.style.transform = `translateY(${rate}px)`;
                    }
                });
            });
        }, { passive: true });
    }
});

// Add CSS custom properties for dynamic theming
document.documentElement.style.setProperty('--nav-transition-duration', '0.25s');
document.documentElement.style.setProperty('--content-transition-duration', '0.3s');

// performance: Preload critical resources
const preloadCriticalResources = () => {
    // Preload fonts if needed
    const fontPreloads = [
        // Add font URLs here if using web fonts
    ];
    
    fontPreloads.forEach(fontUrl => {
        const link = document.createElement('link');
        link.rel = 'preload';
        link.href = fontUrl;
        link.as = 'font';
        link.type = 'font/woff2';
        link.crossOrigin = 'anonymous';
        document.head.appendChild(link);
    });
};

// Initialize performance optimizations
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', preloadCriticalResources);
} else {
    preloadCriticalResources();
}

// Theme toggle functionality
const ThemeManager = {
    // Theme storage key
    STORAGE_KEY: 'streams-to-river-theme',
    
    // Available themes
    themes: {
        LIGHT: 'light',
        DARK: 'dark'
    },
    
    // Initialize theme system
    init() {
        this.themeToggle = document.getElementById('theme-toggle');
        this.currentTheme = this.getSavedTheme() || this.getSystemTheme();
        
        // Apply initial theme
        this.applyTheme(this.currentTheme);
        
        // Setup event listeners
        if (this.themeToggle) {
            this.themeToggle.addEventListener('click', () => this.toggleTheme());
            
            // Add keyboard support
            this.themeToggle.addEventListener('keydown', (e) => {
                if (e.key === 'Enter' || e.key === ' ') {
                    e.preventDefault();
                    this.toggleTheme();
                }
            });
        }
        
        // Listen for system theme changes
        if (window.matchMedia) {
            const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
            mediaQuery.addEventListener('change', (e) => {
                // Only update if user hasn't set a preference
                if (!localStorage.getItem(this.STORAGE_KEY)) {
                    this.applyTheme(e.matches ? this.themes.DARK : this.themes.LIGHT);
                }
            });
        }
    },
    
    // Get saved theme from localStorage
    getSavedTheme() {
        return localStorage.getItem(this.STORAGE_KEY);
    },
    
    // Get system theme preference
    getSystemTheme() {
        if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
            return this.themes.DARK;
        }
        return this.themes.LIGHT;
    },
    
    // Apply theme to document
    applyTheme(theme) {
        this.currentTheme = theme;
        
        // Update document attribute
        document.documentElement.setAttribute('data-theme', theme);
        
        // Update button title
        if (this.themeToggle) {
            const title = theme === this.themes.DARK ? '切换到亮色主题' : '切换到暗色主题';
            this.themeToggle.setAttribute('title', title);
            this.themeToggle.setAttribute('aria-label', title);
        }
        
        // Save to localStorage
        localStorage.setItem(this.STORAGE_KEY, theme);
        
        // Add smooth transition class
        document.documentElement.classList.add('theme-transitioning');
        
        // Remove transition class after animation completes
        setTimeout(() => {
            document.documentElement.classList.remove('theme-transitioning');
        }, 300);
        
        // Trigger custom event for other components
        window.dispatchEvent(new CustomEvent('themeChanged', {
            detail: { theme: theme }
        }));
    },
    
    // Toggle between themes
    toggleTheme() {
        const newTheme = this.currentTheme === this.themes.LIGHT 
            ? this.themes.DARK 
            : this.themes.LIGHT;
        
        // Add click animation
        if (this.themeToggle) {
            this.themeToggle.style.transform = 'scale(0.95)';
            setTimeout(() => {
                this.themeToggle.style.transform = '';
            }, 150);
        }
        
        this.applyTheme(newTheme);
    },
    
    // Get current theme
    getCurrentTheme() {
        return this.currentTheme;
    },
    
    // Check if dark theme is active
    isDarkTheme() {
        return this.currentTheme === this.themes.DARK;
    }
};

// Initialize theme manager when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    ThemeManager.init();
});

// Export for global access
window.ThemeManager = ThemeManager;

// Mermaid integration and management
const MermaidManager = {
    // Initialize Mermaid with theme support
    init() {
        if (typeof mermaid === 'undefined') {
            console.warn('Mermaid library not loaded');
            return;
        }

        // Prevent multiple initializations
        if (window.mermaidInitialized) {
            return;
        }
        window.mermaidInitialized = true;

        console.log('Initializing Mermaid...');

        // Configure Mermaid with initial theme
        this.updateMermaidTheme(ThemeManager.getCurrentTheme());
        
        // Listen for theme changes
        window.addEventListener('themeChanged', (e) => {
            this.updateMermaidTheme(e.detail.theme);
        });

        // Initialize Mermaid
        mermaid.initialize({
            startOnLoad: false,
            securityLevel: 'loose',
            theme: this.getMermaidTheme(ThemeManager.getCurrentTheme()),
            themeVariables: this.getThemeVariables(ThemeManager.getCurrentTheme()),
            fontFamily: 'Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", system-ui, sans-serif',
            flowchart: {
                useMaxWidth: true,
                htmlLabels: true,
                curve: 'basis',
                padding: 20
            },
            sequence: {
                useMaxWidth: true,
                wrap: true,
                padding: 20
            },
            er: {
                useMaxWidth: true,
                padding: 20
            },
            gantt: {
                useMaxWidth: true,
                padding: 20
            },
            journey: {
                useMaxWidth: true,
                padding: 20
            },
            timeline: {
                useMaxWidth: true,
                padding: 20
            },
            graph: {
                useMaxWidth: true,
                padding: 20
            }
        });

        // Render existing diagrams
        this.renderDiagrams();
    },

    // Update Mermaid theme based on current theme
    updateMermaidTheme(theme) {
        if (typeof mermaid === 'undefined') return;

        const mermaidTheme = this.getMermaidTheme(theme);
        const themeVariables = this.getThemeVariables(theme);

        // Reinitialize Mermaid with new theme
        mermaid.initialize({
            startOnLoad: false,
            securityLevel: 'loose',
            theme: mermaidTheme,
            themeVariables: themeVariables,
            fontFamily: 'Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", system-ui, sans-serif',
            flowchart: {
                useMaxWidth: true,
                htmlLabels: true,
                curve: 'basis',
                padding: 20
            },
            sequence: {
                useMaxWidth: true,
                wrap: true,
                padding: 20
            },
            er: {
                useMaxWidth: true,
                padding: 20,
                entityPadding: 18,
                fontSize: 11,
                rowHeight: 25,
                entityRowHeight: 25
            },
            gantt: {
                useMaxWidth: true,
                padding: 20
            },
            journey: {
                useMaxWidth: true,
                padding: 20
            },
            timeline: {
                useMaxWidth: true,
                padding: 20
            },
            graph: {
                useMaxWidth: true,
                padding: 20
            }
        });

        this.renderDiagrams(true); // Force re-render with new theme
    },

    // Get appropriate Mermaid theme based on current theme
    getMermaidTheme(theme) {
        return theme === 'dark' ? 'dark' : 'default';
    },

    // Get theme variables for Mermaid customization
    getThemeVariables(theme) {
        if (theme === 'dark') {
            return {
                // Dark theme variables with updated colors
                primaryColor: '#60a5fa',
                primaryTextColor: '#f7fafc',
                primaryBorderColor: '#3b82f6',
                lineColor: '#e2e8f0',
                secondaryColor: '#2d3748',
                tertiaryColor: '#4a5568',
                background: '#1a202c',
                mainBkg: '#2d3748',
                secondBkg: '#4a5568',
                tertiaryBkg: '#718096',
                nodeBkg: '#2d3748',
                nodeBorder: '#60a5fa',
                clusterBkg: '#4a5568',
                clusterBorder: '#60a5fa',
                defaultLinkColor: '#e2e8f0',
                titleColor: '#f7fafc',
                edgeLabelBackground: '#4a5568',
                nodeTextColor: '#f7fafc',
                fontFamily: 'Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", system-ui, sans-serif',
                // ER diagram specific variables
                attributeBackgroundColorOdd: '#4a5568',
                attributeBackgroundColorEven: '#2d3748',
                entityBackgroundColor: '#2d3748',
                entityBorderColor: '#60a5fa',
                entityTextColor: '#f7fafc',
                erEntityBoxHeight: '22',
                erAttributeBoxHeight: '22',
                // Additional text color overrides for dark theme
                textColor: '#f7fafc',
                cScale0: '#f7fafc',
                cScale1: '#e2e8f0',
                cScale2: '#a0aec0',
                labelTextColor: '#f7fafc',
                sectionBkgColor: '#2d3748',
                altSectionBkgColor: '#4a5568',
                gridColor: '#718096',
                graphDefaultTextColor: '#f7fafc',
                edgeLabelColor: '#e2e8f0'
            };
        } else {
            return {
                // Light theme variables
                primaryColor: '#3b82f6',
                primaryTextColor: '#0f172a',
                primaryBorderColor: '#60a5fa',
                lineColor: '#e2e8f0',
                secondaryColor: '#f8fafc',
                tertiaryColor: '#f1f5f9',
                background: '#ffffff',
                mainBkg: '#ffffff',
                secondBkg: '#f8fafc',
                tertiaryBkg: '#f1f5f9',
                nodeBkg: '#ffffff',
                nodeBorder: '#3b82f6',
                clusterBkg: '#f8fafc',
                clusterBorder: '#3b82f6',
                defaultLinkColor: '#64748b',
                titleColor: '#0f172a',
                edgeLabelBackground: '#ffffff',
                nodeTextColor: '#0f172a',
                fontFamily: 'Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", system-ui, sans-serif',
                // ER diagram specific variables
                attributeBackgroundColorOdd: '#f8fafc',
                attributeBackgroundColorEven: '#ffffff',
                entityBackgroundColor: '#ffffff',
                entityBorderColor: '#3b82f6',
                entityTextColor: '#0f172a',
                erEntityBoxHeight: '22',
                erAttributeBoxHeight: '22'
            };
        }
    },

    // Render all Mermaid diagrams in the document
    async renderDiagrams(forceRerender = false) {
        if (typeof mermaid === 'undefined') return;

        // Find all mermaid code blocks and existing rendered diagrams
        const mermaidBlocks = document.querySelectorAll('pre code.language-mermaid, .mermaid');
        const existingWrappers = document.querySelectorAll('.mermaid-wrapper');
        
        // If forcing re-render, remove existing rendered diagrams
        if (forceRerender) {
            existingWrappers.forEach(wrapper => {
                // Try to find the original content if it exists
                const originalContent = wrapper.getAttribute('data-original-content');
                if (originalContent) {
                    // Recreate the original element structure
                    const pre = document.createElement('pre');
                    const code = document.createElement('code');
                    code.className = 'language-mermaid';
                    code.textContent = originalContent;
                    pre.appendChild(code);
                    wrapper.parentNode.replaceChild(pre, wrapper);
                }
            });
            
            // Re-query after removing existing diagrams
            const newMermaidBlocks = document.querySelectorAll('pre code.language-mermaid, .mermaid');
            for (let i = 0; i < newMermaidBlocks.length; i++) {
                const block = newMermaidBlocks[i];
                await this.renderSingleDiagram(block, i);
            }
        } else {
            // Normal rendering - only render blocks that haven't been rendered yet
            for (let i = 0; i < mermaidBlocks.length; i++) {
                const block = mermaidBlocks[i];
                await this.renderSingleDiagram(block, i);
            }
        }
        
        console.log(`Rendered ${forceRerender ? 'with force re-render' : 'normally'}: found ${mermaidBlocks.length} mermaid blocks`);
    },

    // Render a single Mermaid diagram
    async renderSingleDiagram(element, index) {
        try {
            // Get the mermaid content
            let content = '';
            let targetElement = element;
            
            if (element.tagName === 'CODE' && element.classList.contains('language-mermaid')) {
                content = element.textContent.trim();
                targetElement = element.parentElement; // the <pre> element
            } else if (element.classList.contains('mermaid')) {
                content = element.textContent.trim();
            } else {
                console.log('Element not a mermaid diagram:', element);
                return;
            }

            if (!content) {
                console.log('No content found for mermaid diagram');
                return;
            }

            console.log('Rendering mermaid diagram:', content.substring(0, 50) + '...');

            // Clean up the content
            content = content.replace(/^\s+|\s+$/g, '');

            // Generate unique ID
            const diagramId = `mermaid-diagram-${index}-${Date.now()}`;

            // Validate and render with Mermaid
            try {
                const { svg } = await mermaid.render(diagramId, content);
                
                // Create enhanced container
                const wrapper = document.createElement('div');
                wrapper.className = 'mermaid-wrapper';
                wrapper.setAttribute('data-original-content', content); // Save original content for re-rendering
                
                const container = document.createElement('div');
                container.className = 'mermaid-diagram';
                container.innerHTML = svg;
                
                wrapper.appendChild(container);
                
                // Replace the original element with the rendered diagram
                targetElement.parentNode.replaceChild(wrapper, targetElement);

                // Add click-to-zoom functionality
                this.addZoomFunctionality(wrapper);

                // Clean up SVG attributes for better styling
                const svgElement = container.querySelector('svg');
                if (svgElement) {
                    svgElement.removeAttribute('height');
                    svgElement.style.maxWidth = '100%';
                    svgElement.style.height = 'auto';
                }

                console.log('Successfully rendered mermaid diagram:', diagramId);
            } catch (parseError) {
                throw new Error(`语法错误: ${parseError.message}`);
            }

        } catch (error) {
            console.error('Mermaid rendering error:', error);
            
            // Create error display
            const errorContainer = document.createElement('div');
            errorContainer.className = 'mermaid-error';
            errorContainer.innerHTML = `⚠️ 图表渲染失败: ${error.message}`;
            
            // Replace element with error message
            if (element.parentNode) {
                element.parentNode.replaceChild(errorContainer, element);
            }
        }
    },

    // Add click-to-zoom functionality for diagrams
    addZoomFunctionality(wrapper) {
        const diagram = wrapper.querySelector('.mermaid-diagram');
        if (!diagram) return;

        diagram.style.cursor = 'zoom-in';
        diagram.addEventListener('click', () => {
            this.showZoomedDiagram(diagram);
        });
    },

    // Show diagram in fullscreen zoom overlay
    showZoomedDiagram(diagram) {
        // Create overlay
        const overlay = document.createElement('div');
        overlay.className = 'mermaid-zoom-overlay';
        overlay.innerHTML = `
            <div class="mermaid-zoom-container">
                <button class="mermaid-zoom-close" title="关闭">×</button>
                <div class="mermaid-zoom-content">
                    ${diagram.innerHTML}
                </div>
            </div>
        `;

        // Add to document
        document.body.appendChild(overlay);

        // Add event listeners
        const closeBtn = overlay.querySelector('.mermaid-zoom-close');
        closeBtn.addEventListener('click', () => {
            document.body.removeChild(overlay);
        });

        overlay.addEventListener('click', (e) => {
            if (e.target === overlay) {
                document.body.removeChild(overlay);
            }
        });

        // ESC key to close
        const handleEsc = (e) => {
            if (e.key === 'Escape') {
                if (document.body.contains(overlay)) {
                    document.body.removeChild(overlay);
                }
                document.removeEventListener('keydown', handleEsc);
            }
        };
        document.addEventListener('keydown', handleEsc);

        // Animate in
        setTimeout(() => {
            overlay.classList.add('active');
        }, 10);
    }
};

// Initialize Mermaid when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    // Wait for Mermaid library to load
    const initMermaid = () => {
        if (typeof mermaid !== 'undefined') {
            console.log('Mermaid library loaded, initializing...');
            MermaidManager.init();
        } else {
            console.log('Waiting for Mermaid library...');
            // If mermaid is not loaded yet, try again
            setTimeout(initMermaid, 200);
        }
    };
    
    // Start initialization
    setTimeout(initMermaid, 100);
});

// Also try to initialize when window loads (backup)
window.addEventListener('load', () => {
    if (typeof mermaid !== 'undefined' && !window.mermaidInitialized) {
        console.log('Window loaded, trying Mermaid init backup...');
        MermaidManager.init();
    }
});

// Export for global access
window.MermaidManager = MermaidManager;

// Sticky Header Manager
const StickyHeaderManager = {
    // Initialize sticky header functionality
    init() {
        this.header = document.querySelector('.content-header');
        this.contentBody = document.querySelector('.content-body');
        this.lastScrollY = 0;
        this.isSticky = false;
        
        if (!this.header || !this.contentBody) return;
        
        // Add sticky behavior on content scroll
        this.contentBody.addEventListener('scroll', () => this.handleScroll(), { passive: true });
        
        // Initialize header state
        this.updateHeaderState();
    },
    
    // Handle scroll events
    handleScroll() {
        if (!this.contentBody || !this.header) return;
        
        const scrollY = this.contentBody.scrollTop;
        const scrollingDown = scrollY > this.lastScrollY;
        const shouldBeSticky = scrollY > 20;
        
        // Update sticky state
        if (shouldBeSticky !== this.isSticky) {
            this.isSticky = shouldBeSticky;
            this.updateHeaderState();
        }
        
        // Handle scroll direction for enhanced effects
        if (Math.abs(scrollY - this.lastScrollY) > 5) {
            this.header.classList.toggle('scrolling-down', scrollingDown && this.isSticky);
            this.header.classList.toggle('scrolling-up', !scrollingDown && this.isSticky);
        }
        
        this.lastScrollY = scrollY;
    },
    
    // Update header visual state
    updateHeaderState() {
        if (!this.header) return;
        
        if (this.isSticky) {
            this.header.classList.add('sticky-active');
        } else {
            this.header.classList.remove('sticky-active');
        }
    }
};

// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    StickyHeaderManager.init();
});

// Export for global access
window.StickyHeaderManager = StickyHeaderManager;

// Code Highlighting Manager
const CodeHighlightManager = {
    // Initialize code highlighting
    init() {
        console.log('Initializing code highlighting...');
        this.highlightAllCodeBlocks();
    },

    // Highlight all code blocks in the document
    highlightAllCodeBlocks() {
        // Try multiple selectors to find code blocks, including indented ones
        const selectors = [
            'pre code',                    // Standard structure
            '.markdown-content pre code',  // Within markdown container
            'pre',                         // Standalone pre elements
            '.markdown-content pre',       // Pre elements in markdown
            'li pre code',                 // Code blocks in list items
            'li pre',                      // Pre elements in list items
            'ol pre code',                 // Code blocks in ordered lists
            'ol pre',                      // Pre elements in ordered lists
            'ul pre code',                 // Code blocks in unordered lists
            'ul pre'                       // Pre elements in unordered lists
        ];
        
        let allCodeBlocks = new Set(); // Use Set to avoid duplicates
        
        for (const selector of selectors) {
            const blocks = document.querySelectorAll(selector);
            blocks.forEach(block => allCodeBlocks.add(block));
        }
        
        const codeBlocks = Array.from(allCodeBlocks);
        console.log(`Found ${codeBlocks.length} total unique code blocks`);
        
        if (codeBlocks.length === 0) {
            console.log('No code blocks found with any selector');
            return;
        }
        
        codeBlocks.forEach((block, index) => {
            this.highlightCodeBlock(block, index);
        });
    },

    // Highlight a single code block
    highlightCodeBlock(codeBlock, index) {
        // Handle both <pre><code> and standalone <pre> elements
        let targetElement = codeBlock;
        let parentPre = null;
        
        if (codeBlock.tagName === 'PRE') {
            // If it's a <pre> element, look for child <code>
            const codeChild = codeBlock.querySelector('code');
            if (codeChild) {
                targetElement = codeChild;
                parentPre = codeBlock;
            } else {
                targetElement = codeBlock;
            }
        } else if (codeBlock.tagName === 'CODE') {
            // If it's a <code> element, get parent <pre>
            parentPre = codeBlock.closest('pre');
            targetElement = codeBlock;
        }

        // Get the language from class name (check both code and pre elements)
        let language = 'text';
        const elementsToCheck = [targetElement, parentPre].filter(el => el);
        
        for (const element of elementsToCheck) {
            if (!element) continue;
            const classList = Array.from(element.classList);
            const langClass = classList.find(cls => cls.startsWith('language-'));
            if (langClass) {
                language = langClass.replace('language-', '');
                break;
            }
        }
        
        // If no language class found, try to detect from content
        if (language === 'text') {
            const content = targetElement.textContent.trim();
            language = this.detectLanguage(content);
        }

        console.log(`Highlighting code block ${index} with language: ${language}, element:`, targetElement.tagName);

        // Apply syntax highlighting based on language
        switch (language) {
            case 'json':
                this.highlightJSON(targetElement);
                break;
            case 'javascript':
            case 'js':
                this.highlightJavaScript(targetElement);
                break;
            case 'go':
                this.highlightGo(targetElement);
                break;
            case 'python':
                this.highlightPython(targetElement);
                break;
            case 'bash':
            case 'shell':
                this.highlightBash(targetElement);
                break;
            case 'curl':
                this.highlightCurl(targetElement);
                break;
            case 'protobuf':
                this.highlightProtobuf(targetElement);
                break;
            case 'sql':
                this.highlightSQL(targetElement);
                break;
            default:
                this.highlightGeneric(targetElement);
                break;
        }

        // Add copy button to the pre element
        const preElement = parentPre || (targetElement.tagName === 'PRE' ? targetElement : targetElement.closest('pre'));
        if (preElement) {
            this.addCopyButton(preElement);
        }
    },

    // Auto-detect language from content
    detectLanguage(content) {
        // Remove leading/trailing whitespace and normalize indentation
        content = content.trim();
        
        // JSON detection - be more flexible with whitespace
        const trimmedContent = content.replace(/^\s+/gm, '').trim();
        if ((trimmedContent.startsWith('{') && trimmedContent.endsWith('}')) || 
            (trimmedContent.startsWith('[') && trimmedContent.endsWith(']'))) {
            try {
                JSON.parse(trimmedContent);
                console.log('Detected JSON content:', trimmedContent.substring(0, 50) + '...');
                return 'json';
            } catch (e) {
                // Try the original content too
                try {
                    JSON.parse(content);
                    console.log('Detected JSON content (original):', content.substring(0, 50) + '...');
                    return 'json';
                } catch (e2) {
                    // Not valid JSON, continue with other checks
                }
            }
        }
        
        // Go detection
        if (content.includes('package ') && content.includes('func ')) {
            return 'go';
        }
        
        // Python detection
        if (content.includes('def ') || content.includes('import ') || content.includes('from ')) {
            return 'python';
        }
        
        // JavaScript detection
        if (content.includes('function ') || content.includes('const ') || content.includes('let ')) {
            return 'javascript';
        }
        
        // Bash detection
        if (content.includes('#!/bin/bash') || content.includes('curl ') || content.startsWith('$ ')) {
            return 'bash';
        }
        
        // SQL detection
        if (/\b(SELECT|INSERT|UPDATE|DELETE|CREATE|ALTER|DROP)\b/i.test(content)) {
            return 'sql';
        }
        
        return 'text';
    },

    // Escape HTML entities
    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    },

    // Highlight JSON
    highlightJSON(codeBlock) {
        console.log('Highlighting JSON block:', codeBlock);
        
        let content = codeBlock.textContent;
        console.log('JSON content:', content.substring(0, 100) + '...');
        
        try {
            // Parse and re-stringify for validation and formatting
            const parsed = JSON.parse(content);
            content = JSON.stringify(parsed, null, 2);
        } catch (e) {
            // If JSON is invalid, keep original content
            console.log('JSON parse failed, using original content');
        }
        
        // Force styles immediately
        codeBlock.style.backgroundColor = '#282c34';
        codeBlock.style.color = '#abb2bf';
        
        // Also style parent pre element
        const parentPre = codeBlock.closest('pre');
        if (parentPre) {
            parentPre.style.backgroundColor = '#282c34';
            parentPre.classList.add('highlighted-json');
            console.log('Styled parent pre element');
        }
        
        // Escape HTML first
        const escaped = this.escapeHtml(content);
        
        // Apply JSON syntax highlighting
        const highlighted = escaped
            .replace(/(&quot;[^&]*?&quot;)\s*:/g, '<span class="json-key">$1</span>:')
            .replace(/:\s*(&quot;[^&]*?&quot;)/g, ': <span class="json-string">$1</span>')
            .replace(/:\s*(\d+\.?\d*)/g, ': <span class="json-number">$1</span>')
            .replace(/:\s*(true|false)/g, ': <span class="json-boolean">$1</span>')
            .replace(/:\s*(null)/g, ': <span class="json-null">$1</span>')
            .replace(/([{}[\],])/g, '<span class="json-punctuation">$1</span>');
        
        codeBlock.innerHTML = highlighted;
        console.log('JSON highlighting applied successfully');
    },

    // Highlight JavaScript
    highlightJavaScript(codeBlock) {
        const content = codeBlock.textContent;
        const escaped = this.escapeHtml(content);
        
        const highlighted = escaped
            // Keywords
            .replace(/\b(const|let|var|function|if|else|for|while|return|class|import|export|from|async|await|try|catch|finally)\b/g, '<span class="js-keyword">$1</span>')
            // Strings (handle escaped quotes)
            .replace(/(&quot;[^&]*?&quot;|&#x27;[^&]*?&#x27;|`[^`]*?`)/g, '<span class="js-string">$1</span>')
            // Numbers
            .replace(/\b(\d+\.?\d*)\b/g, '<span class="js-number">$1</span>')
            // Comments
            .replace(/(\/\/[^\n]*)/g, '<span class="js-comment">$1</span>')
            .replace(/(\/\*[\s\S]*?\*\/)/g, '<span class="js-comment">$1</span>');
        
        codeBlock.innerHTML = highlighted;
    },

    // Highlight Go
    highlightGo(codeBlock) {
        const content = codeBlock.textContent;
        const escaped = this.escapeHtml(content);
        
        const highlighted = escaped
            // Keywords
            .replace(/\b(package|import|func|var|const|type|struct|interface|if|else|for|range|return|go|defer|chan|select|case|default|switch|break|continue)\b/g, '<span class="go-keyword">$1</span>')
            // Types
            .replace(/\b(string|int|int32|int64|float32|float64|bool|byte|rune|error)\b/g, '<span class="go-type">$1</span>')
            // Strings
            .replace(/(`[^`]*`|&quot;[^&]*?&quot;)/g, '<span class="go-string">$1</span>')
            // Numbers
            .replace(/\b(\d+\.?\d*)\b/g, '<span class="go-number">$1</span>')
            // Comments
            .replace(/(\/\/[^\n]*)/g, '<span class="go-comment">$1</span>');
        
        codeBlock.innerHTML = highlighted;
    },

    // Highlight Python
    highlightPython(codeBlock) {
        const content = codeBlock.textContent;
        const escaped = this.escapeHtml(content);
        
        const highlighted = escaped
            // Keywords
            .replace(/\b(def|class|if|elif|else|for|while|return|import|from|as|try|except|finally|with|lambda|and|or|not|in|is|None|True|False)\b/g, '<span class="python-keyword">$1</span>')
            // Strings
            .replace(/(&quot;&quot;&quot;[\s\S]*?&quot;&quot;&quot;|&#x27;&#x27;&#x27;[\s\S]*?&#x27;&#x27;&#x27;|&quot;[^&]*?&quot;|&#x27;[^&]*?&#x27;)/g, '<span class="python-string">$1</span>')
            // Numbers
            .replace(/\b(\d+\.?\d*)\b/g, '<span class="python-number">$1</span>')
            // Comments
            .replace(/(#[^\n]*)/g, '<span class="python-comment">$1</span>');
        
        codeBlock.innerHTML = highlighted;
    },

    // Highlight Bash/Shell
    highlightBash(codeBlock) {
        const content = codeBlock.textContent;
        const escaped = this.escapeHtml(content);
        
        const highlighted = escaped
            // Commands
            .replace(/\b(curl|wget|ls|cd|mkdir|rm|cp|mv|grep|awk|sed|cat|echo|export|source)\b/g, '<span class="bash-command">$1</span>')
            // Options
            .replace(/(\s-{1,2}[a-zA-Z-]+)/g, '<span class="bash-option">$1</span>')
            // Strings
            .replace(/(&quot;[^&]*?&quot;|&#x27;[^&]*?&#x27;)/g, '<span class="bash-string">$1</span>')
            // Comments
            .replace(/(#[^\n]*)/g, '<span class="bash-comment">$1</span>');
        
        codeBlock.innerHTML = highlighted;
    },

    // Highlight Curl commands
    highlightCurl(codeBlock) {
        const content = codeBlock.textContent;
        const escaped = this.escapeHtml(content);
        
        const highlighted = escaped
            // curl command
            .replace(/\bcurl\b/g, '<span class="curl-command">curl</span>')
            // HTTP methods
            .replace(/(-X\s+)(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)/g, '$1<span class="curl-method">$2</span>')
            // Headers
            .replace(/(-H\s+)(&quot;[^&]*?&quot;)/g, '$1<span class="curl-header">$2</span>')
            // URLs
            .replace(/(https?:\/\/[^\s&]*)/g, '<span class="curl-url">$1</span>')
            // Options
            .replace(/(\s-{1,2}[a-zA-Z-]+)/g, '<span class="curl-option">$1</span>');
        
        codeBlock.innerHTML = highlighted;
    },

    // Highlight Protobuf
    highlightProtobuf(codeBlock) {
        const content = codeBlock.textContent;
        const escaped = this.escapeHtml(content);
        
        const highlighted = escaped
            // Keywords
            .replace(/\b(message|service|rpc|returns|repeated|optional|required|enum|import|package|option|syntax)\b/g, '<span class="proto-keyword">$1</span>')
            // Types
            .replace(/\b(string|int32|int64|uint32|uint64|bool|bytes|double|float)\b/g, '<span class="proto-type">$1</span>')
            // Field numbers
            .replace(/=\s*(\d+)/g, '= <span class="proto-number">$1</span>')
            // Strings
            .replace(/(&quot;[^&]*?&quot;)/g, '<span class="proto-string">$1</span>')
            // Comments
            .replace(/(\/\/[^\n]*)/g, '<span class="proto-comment">$1</span>');
        
        codeBlock.innerHTML = highlighted;
    },

    // Highlight SQL
    highlightSQL(codeBlock) {
        const content = codeBlock.textContent;
        const escaped = this.escapeHtml(content);
        
        const highlighted = escaped
            // Keywords
            .replace(/\b(SELECT|FROM|WHERE|INSERT|UPDATE|DELETE|CREATE|ALTER|DROP|TABLE|INDEX|PRIMARY|KEY|FOREIGN|REFERENCES|NOT|NULL|DEFAULT|AUTO_INCREMENT|UNIQUE|AND|OR|ORDER|BY|GROUP|HAVING|LIMIT|OFFSET|JOIN|LEFT|RIGHT|INNER|OUTER|ON|AS|DISTINCT|COUNT|SUM|AVG|MAX|MIN)\b/gi, '<span class="sql-keyword">$1</span>')
            // Strings
            .replace(/(&#x27;[^&]*?&#x27;)/g, '<span class="sql-string">$1</span>')
            // Numbers
            .replace(/\b(\d+\.?\d*)\b/g, '<span class="sql-number">$1</span>')
            // Comments
            .replace(/(--[^\n]*)/g, '<span class="sql-comment">$1</span>');
        
        codeBlock.innerHTML = highlighted;
    },

    // Generic highlighting for unknown languages
    highlightGeneric(codeBlock) {
        // Just ensure the content is properly escaped
        const content = codeBlock.textContent;
        codeBlock.textContent = content; // This will escape HTML entities
    },

    // Add copy button to code block
    addCopyButton(codeBlock) {
        const pre = codeBlock.parentElement;
        if (!pre || pre.tagName !== 'PRE') return;

        // Check if copy button already exists
        if (pre.querySelector('.code-copy-btn')) return;

        const copyBtn = document.createElement('button');
        copyBtn.className = 'code-copy-btn';
        copyBtn.innerHTML = '📋';
        copyBtn.title = '复制代码';
        copyBtn.setAttribute('aria-label', '复制代码');

        copyBtn.addEventListener('click', async () => {
            try {
                // Get the original text content without HTML tags
                const text = codeBlock.textContent || codeBlock.innerText;
                await navigator.clipboard.writeText(text);
                
                // Show feedback
                copyBtn.innerHTML = '✅';
                copyBtn.title = '已复制';
                
                setTimeout(() => {
                    copyBtn.innerHTML = '📋';
                    copyBtn.title = '复制代码';
                }, 2000);
            } catch (err) {
                console.error('Failed to copy code:', err);
                copyBtn.innerHTML = '❌';
                setTimeout(() => {
                    copyBtn.innerHTML = '📋';
                }, 2000);
            }
        });

        pre.style.position = 'relative';
        pre.appendChild(copyBtn);
    }
};

// Initialize Code Highlighting - multiple triggers to ensure it runs
function initCodeHighlighting() {
    console.log('Attempting to initialize code highlighting...');
    CodeHighlightManager.init();
}

// Try multiple initialization triggers
document.addEventListener('DOMContentLoaded', () => {
    console.log('DOMContentLoaded - initializing code highlighting');
    setTimeout(initCodeHighlighting, 100);
    setTimeout(initCodeHighlighting, 500);
    setTimeout(initCodeHighlighting, 1000);
});

// Also try when window loads (in case DOMContentLoaded already fired)
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initCodeHighlighting);
} else {
    // DOM already loaded
    console.log('DOM already loaded - initializing immediately');
    initCodeHighlighting();
    setTimeout(initCodeHighlighting, 100);
}

// Also reinitialize when content changes
const observer = new MutationObserver((mutations) => {
    let shouldReinit = false;
    mutations.forEach((mutation) => {
        if (mutation.type === 'childList') {
            // Check if new code blocks were added
            const addedNodes = Array.from(mutation.addedNodes);
            if (addedNodes.some(node => 
                node.nodeType === 1 && 
                (node.tagName === 'PRE' || node.querySelector('pre'))
            )) {
                shouldReinit = true;
            }
        }
    });
    
    if (shouldReinit) {
        setTimeout(() => {
            CodeHighlightManager.init();
        }, 100);
    }
});

// Start observing document changes
observer.observe(document.body, {
    childList: true,
    subtree: true
});

// Export for global access
window.CodeHighlightManager = CodeHighlightManager;

// ============ 置顶按钮管理器 ============
const BackToTopManager = {
    init() {
        this.button = document.getElementById('back-to-top');
        if (!this.button) {
            console.warn('Back to top button not found');
            return;
        }

        this.setupEventListeners();
        this.handleScroll(); // Check initial state
    },

    setupEventListeners() {
        // Add click event listener for smooth scroll
        this.button.addEventListener('click', (e) => {
            e.preventDefault();
            this.scrollToTop();
        });

        // Add scroll event listener with throttling
        let scrollTimeout;
        window.addEventListener('scroll', () => {
            if (scrollTimeout) {
                clearTimeout(scrollTimeout);
            }
            scrollTimeout = setTimeout(() => {
                this.handleScroll();
            }, 16); // ~60fps
        });

        // Add keyboard support
        this.button.addEventListener('keydown', (e) => {
            if (e.key === 'Enter' || e.key === ' ') {
                e.preventDefault();
                this.scrollToTop();
            }
        });
    },

    handleScroll() {
        const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
        const threshold = 300; // Show button after 300px scroll

        if (scrollTop > threshold) {
            this.showButton();
        } else {
            this.hideButton();
        }
    },

    showButton() {
        if (!this.button.classList.contains('visible')) {
            this.button.classList.add('visible');
            this.button.setAttribute('aria-hidden', 'false');
        }
    },

    hideButton() {
        if (this.button.classList.contains('visible')) {
            this.button.classList.remove('visible');
            this.button.setAttribute('aria-hidden', 'true');
        }
    },

    scrollToTop() {
        // Use the existing smooth scroll function if available
        if (window.smoothScrollTo) {
            window.smoothScrollTo(0, 800);
        } else {
            // Fallback smooth scroll implementation
            const startTime = performance.now();
            const startPosition = window.pageYOffset;
            const duration = 800;

            const easeOutExpo = (t) => {
                return t === 1 ? 1 : 1 - Math.pow(2, -10 * t);
            };

            const animateScroll = (currentTime) => {
                const elapsed = currentTime - startTime;
                const progress = Math.min(elapsed / duration, 1);
                const easedProgress = easeOutExpo(progress);
                
                const currentPosition = startPosition * (1 - easedProgress);
                window.scrollTo(0, currentPosition);

                if (progress < 1) {
                    requestAnimationFrame(animateScroll);
                }
            };

            requestAnimationFrame(animateScroll);
        }

        // Add visual feedback
        this.button.style.transform = 'translateY(-1px) scale(0.95)';
        setTimeout(() => {
            this.button.style.transform = '';
        }, 150);
    }
};

// Initialize back to top functionality
document.addEventListener('DOMContentLoaded', () => {
    console.log('Initializing Back to Top Manager...');
    BackToTopManager.init();
});

// Also initialize if DOM is already loaded
if (document.readyState !== 'loading') {
    BackToTopManager.init();
}

// Export for global access
window.BackToTopManager = BackToTopManager;
