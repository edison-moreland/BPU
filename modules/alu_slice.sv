module alu_slice (
    input  logic a, b,
    input  logic prev_eq, prev_greater,
    output logic result_xor,
    output logic result_or,
    output logic result_and,
    output logic result_not,
    output logic eq, greater,
);

    always @* begin
        eq <= (a == b) & prev_eq;
        greater <= ((a > b) & prev_eq) | prev_greater;

        result_xor <= a ^ b;
        result_or <= a | b;
        result_and <= a & b;
        result_not <= !a;
    end

endmodule