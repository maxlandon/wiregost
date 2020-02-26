#include <windows.h>

void RunGhost();

BOOL WINAPI DllMain(
    HINSTANCE _hinstDLL, // handle to DLL module
    DWORD _fdwReason,    // reason for calling function
    LPVOID _lpReserved   // reserved
);
