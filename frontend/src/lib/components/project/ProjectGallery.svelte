<script lang="ts">
    import type { ImgType, ImageSrc } from '$lib/types';
    import Gallery from "../core/gallery/Gallery.svelte";

    const transformImage = (src: string): ImageSrc => ({
        img: { src, w: 800, h: 800 }, // Assuming width and height based on your URLs
        sources: {
            // You can add different formats and sizes here if available
            //jpeg: [{ src, w: 800, h: 800 }],
        },
    });
    
    const images1: ImgType[] = [
        { alt: 'example', src: transformImage('https://images.unsplash.com/photo-1687851898206-6636e30dcab1?ixid=M3w0Njc5ODF8MHwxfGFsbHx8fHx8fHx8fDE2ODc5NzExMDB8&ixlib=rb-4.0.3&w=800&h=800&auto=format&fit=max') },
        { alt: 'example', src: transformImage('https://images.unsplash.com/photo-1687333517042-b1749fc15642?ixid=M3w0Njc5ODF8MHwxfGFsbHx8fHx8fHx8fDE2ODc5NzY2MzZ8&ixlib=rb-4.0.3&w=800&h=800&auto=format&fit=max') },
        { alt: 'example', src: transformImage('https://images.unsplash.com/photo-1676149841755-12f0b913737c?ixid=M3w0Njc5ODF8MHwxfGFsbHx8fHx8fHx8fDE2ODc5NzY3MzR8&ixlib=rb-4.0.3&w=800&h=800&auto=format&fit=max') }
    ];

    const images2: ImgType[] = [
        { alt: 'example', src: transformImage('https://images.unsplash.com/photo-1687593883595-6b678fc5c413?ixid=M3w0Njc5ODF8MHwxfGFsbHx8fHx8fHx8fDE2ODc5NzY3NjF8&ixlib=rb-4.0.3&w=800&h=800&auto=format&fit=max') },
        { alt: 'example', src: transformImage('https://images.unsplash.com/photo-1652959894563-e815de31350e?ixid=M3w0Njc5ODF8MHwxfGFsbHx8fHx8fHx8fDE2ODc5NzY3NzN8&ixlib=rb-4.0.3&w=800&h=800&auto=format&fit=max') },
        { alt: 'example', src: transformImage('https://images.unsplash.com/photo-1632073594175-a9278f2ebb8b?ixid=M3w0Njc5ODF8MHwxfGFsbHx8fHx8fHx8fDE2ODc5NzY2NTF8&ixlib=rb-4.0.3&w=800&h=800&auto=format&fit=max') }
    ];

    const images3: ImgType[] = [
        { alt: 'example', src: transformImage('https://images.unsplash.com/photo-1681845719580-8fd75eced4d7?ixid=M3w0Njc5ODF8MHwxfGFsbHx8fHx8fHx8fDE2ODc5NzY3ODJ8&ixlib=rb-4.0.3&w=800&h=800&auto=format&fit=max') },
        { alt: 'example', src: transformImage('https://images.unsplash.com/photo-1631465416799-02880b27977e?ixid=M3w0Njc5ODF8MHwxfGFsbHx8fHx8fHx8fDE2ODc5NzY2NjF8&ixlib=rb-4.0.3&w=800&h=800&auto=format&fit=max') },
        { alt: 'example', src: transformImage('https://images.unsplash.com/photo-1679828485992-95fb004ec4f9?ixid=M3w0Njc5ODF8MHwxfGFsbHx8fHx8fHx8fDE2ODc5NzY3OTV8&ixlib=rb-4.0.3&w=800&h=800&auto=format&fit=max') }
    ];

    const images4: ImgType[] = [
        { alt: 'example', src: transformImage('https://images.unsplash.com/photo-1626553684728-ac5f563eefe7?ixid=M3w0Njc5ODF8MHwxfGFsbHx8fHx8fHx8fDE2ODc5NzY4MDV8&ixlib=rb-4.0.3&w=800&h=800&auto=format&fit=max') },
        { alt: 'example', src: transformImage('https://images.unsplash.com/photo-1680456265112-e4115432ef23?ixid=M3w0Njc5ODF8MHwxfGFsbHx8fHx8fHx8fDE2ODc5NzY4MTZ8&ixlib=rb-4.0.3&w=800&h=800&auto=format&fit=max') },
        { alt: 'example', src: transformImage('https://images.unsplash.com/photo-1632687222219-93f0ecb7dab0?ixid=M3w0Njc5ODF8MHwxfGFsbHx8fHx8fHx8fDE2ODc5NzY4NjV8&ixlib=rb-4.0.3&w=800&h=800&auto=format&fit=max') }
    ];
</script>

<Gallery class="gap-2 grid-cols-2 md:grid-cols-4">
    <Gallery items={images1} />
    <Gallery items={images2} />
    <Gallery items={images3} />
    <Gallery items={images4} />
</Gallery>