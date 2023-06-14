install(
    TARGETS gpio_exe
    RUNTIME COMPONENT gpio_Runtime
)

if(PROJECT_IS_TOP_LEVEL)
  include(CPack)
endif()
